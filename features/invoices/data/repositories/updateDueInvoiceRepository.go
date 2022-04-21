package invoicedomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateDueInvoiceRepository struct {
	invoiceDataSource                 databaseinvoicedatasourceinterfaces.InvoiceDataSource
	memberAccessDataSource            databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	updateInvoiceTransactionComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent
	createNotificationComponent       notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction                mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                      string
}

func NewUpdateDueInvoiceRepository(
	invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	updateDueInvoiceRepositoryTransactionComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (invoicedomainrepositoryinterfaces.UpdateDueInvoiceRepository, error) {
	updateDueInvoiceRepo := &updateDueInvoiceRepository{
		invoiceDataSource,
		memberAccessDataSource,
		updateDueInvoiceRepositoryTransactionComponent,
		createNotificationComponent,
		mongoDBTransaction,
		"UpdateDueInvoiceRepository",
	}

	mongoDBTransaction.SetTransaction(
		updateDueInvoiceRepo,
		"UpdateDueInvoiceRepository",
	)

	return updateDueInvoiceRepo, nil
}

func (updateDueInvoiceRepo *updateDueInvoiceRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateDueInvoiceRepo.updateInvoiceTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateInvoice),
	)
}

func (updateDueInvoiceRepo *updateDueInvoiceRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	invoiceToUpdate := input.(*model.InternalUpdateInvoice)

	return updateDueInvoiceRepo.updateInvoiceTransactionComponent.TransactionBody(
		operationOption,
		invoiceToUpdate,
	)
}

func (updateDueInvoiceRepo *updateDueInvoiceRepository) RunTransaction() ([]*model.Invoice, error) {
	loc, _ := time.LoadLocation("Asia/Bangkok")

	invoices, err := updateDueInvoiceRepo.invoiceDataSource.GetMongoDataSource().Find(
		map[string]interface{}{
			"paymentDueDate": map[string]interface{}{
				"$lte": time.Now().In(loc),
			},
			"status": model.InvoiceStatusAvailable,
		},
		&mongodbcoretypes.PaginationOptions{},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateDueInvoiceRepo.pathIdentity,
			err,
		)
	}

	for _, invoice := range invoices {
		output, err := updateDueInvoiceRepo.mongoDBTransaction.RunTransaction(
			&model.InternalUpdateInvoice{
				ID: invoice.ID,
				Status: func(s model.InvoiceStatus) *model.InvoiceStatus {
					return &s
				}(model.InvoiceStatusPaymentNeeded),
			},
		)
		if err != nil {
			return nil, err
		}

		updatedInvoice := (output).(*model.Invoice)
		go func() {
			memberAccessesToNotify, err := updateDueInvoiceRepo.memberAccessDataSource.GetMongoDataSource().Find(
				map[string]interface{}{
					"organization._id":   updatedInvoice.Organization.ID,
					"status":             model.MemberAccessStatusActive,
					"proposalStatus":     model.EntityProposalStatusApproved,
					"invitationAccepted": true,
				},
				&mongodbcoretypes.PaginationOptions{},
				&mongodbcoretypes.OperationOptions{},
			)
			if err != nil {
				return
			}

			jsonInvPayload, _ := json.Marshal(updatedInvoice)
			for _, memberAccess := range memberAccessesToNotify {
				notifToCreate := &model.InternalCreateNotification{
					NotificationCategory: model.NotificationCategoryInvoiceUpdatedPaymentNeeded,
					RecipientAccount: &model.ObjectIDOnly{
						ID: &memberAccess.Account.ID,
					},
					PayloadOptions: &model.PayloadOptionsInput{
						InvoicePayload: &model.InvoicePayloadInput{
							Invoice: &model.InvoiceForNotifPayloadInput{},
						},
					},
				}
				json.Unmarshal(jsonInvPayload, &notifToCreate.PayloadOptions.InvoicePayload.Invoice)
				_, err := updateDueInvoiceRepo.createNotificationComponent.TransactionBody(
					&mongodbcoretypes.OperationOptions{},
					notifToCreate,
				)
				if err != nil {
					return
				}
			}
		}()
	}

	return invoices, nil
}
