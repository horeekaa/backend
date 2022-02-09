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

type createInvoiceRepository struct {
	memberAccessDataSource            databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	createInvoiceTransactionComponent invoicedomainrepositoryinterfaces.CreateInvoiceTransactionComponent
	createNotificationComponent       notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction                mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateInvoiceRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	createInvoiceRepositoryTransactionComponent invoicedomainrepositoryinterfaces.CreateInvoiceTransactionComponent,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (invoicedomainrepositoryinterfaces.CreateInvoiceRepository, error) {
	createInvoiceRepo := &createInvoiceRepository{
		memberAccessDataSource,
		createInvoiceRepositoryTransactionComponent,
		createNotificationComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createInvoiceRepo,
		"CreateInvoiceRepository",
	)

	return createInvoiceRepo, nil
}

func (createInvoiceRepo *createInvoiceRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createInvoiceRepo.createInvoiceTransactionComponent.PreTransaction(
		input.(*model.InternalCreateInvoice),
	)
}

func (createInvoiceRepo *createInvoiceRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	invoiceToCreate := input.(*model.InternalCreateInvoice)

	return createInvoiceRepo.createInvoiceTransactionComponent.TransactionBody(
		operationOption,
		invoiceToCreate,
	)
}

func (createInvoiceRepo *createInvoiceRepository) RunTransaction(
	input *model.InternalCreateInvoice,
) ([]*model.Invoice, error) {
	output, err := createInvoiceRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	createdInvoices := (output).([]*model.Invoice)
	go func() {
		for _, invoice := range createdInvoices {
			memberAccessesToNotify, err := createInvoiceRepo.memberAccessDataSource.GetMongoDataSource().Find(
				map[string]interface{}{
					"organization._id":   invoice.Organization.ID,
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

			jsonInvPayload, _ := json.Marshal(invoice)
			for _, memberAccess := range memberAccessesToNotify {
				notifToCreate := &model.InternalCreateNotification{
					NotificationCategory: model.NotificationCategoryInvoiceCreated,
					RecipientAccount: &model.ObjectIDOnly{
						ID: &memberAccess.Account.ID,
					},
					PayloadOptions: &model.PayloadOptionsInput{
						InvoiceCreatedPayload: &model.InvoiceCreatedPayloadInput{
							Invoice: &model.InvoiceForNotifPayloadInput{},
						},
					},
				}
				json.Unmarshal(jsonInvPayload, &notifToCreate.PayloadOptions.InvoiceCreatedPayload.Invoice)
				_, err := createInvoiceRepo.createNotificationComponent.TransactionBody(
					&mongodbcoretypes.OperationOptions{},
					notifToCreate,
				)
				if err != nil {
					return
				}
			}
		}
	}()
	return createdInvoices, nil
}
