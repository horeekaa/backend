package invoicedomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateInvoiceRepository struct {
	memberAccessDataSource            databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	updateInvoiceTransactionComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent
	createNotificationComponent       notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction                mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewUpdateInvoiceRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	updateInvoiceRepositoryTransactionComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (invoicedomainrepositoryinterfaces.UpdateInvoiceRepository, error) {
	updateInvoiceRepo := &updateInvoiceRepository{
		memberAccessDataSource,
		updateInvoiceRepositoryTransactionComponent,
		createNotificationComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		updateInvoiceRepo,
		"updateInvoiceRepository",
	)

	return updateInvoiceRepo, nil
}

func (updateInvoiceRepo *updateInvoiceRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateInvoiceRepo.updateInvoiceTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateInvoice),
	)
}

func (updateInvoiceRepo *updateInvoiceRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	invoiceToUpdate := input.(*model.InternalUpdateInvoice)

	return updateInvoiceRepo.updateInvoiceTransactionComponent.TransactionBody(
		operationOption,
		invoiceToUpdate,
	)
}

func (updateInvoiceRepo *updateInvoiceRepository) RunTransaction(
	input *model.InternalUpdateInvoice,
) (*model.Invoice, error) {
	output, err := updateInvoiceRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	updatedInvoice := (output).(*model.Invoice)
	go func() {
		memberAccessesToNotify, err := updateInvoiceRepo.memberAccessDataSource.GetMongoDataSource().Find(
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
				NotificationCategory: model.NotificationCategoryInvoiceUpdatedPlain,
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
			_, err := updateInvoiceRepo.createNotificationComponent.TransactionBody(
				&mongodbcoretypes.OperationOptions{},
				notifToCreate,
			)
			if err != nil {
				return
			}
		}
	}()
	return updatedInvoice, nil
}
