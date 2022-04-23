package invoicedomainrepositories

import (
	"encoding/json"
	"reflect"
	"time"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	invoicedomainrepositorytypes "github.com/horeekaa/backend/features/invoices/domain/repositories/types"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	"github.com/horeekaa/backend/model"
)

type createInvoiceRepository struct {
	purchaseOrderDataSource           databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	memberAccessDataSource            databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	createInvoiceTransactionComponent invoicedomainrepositoryinterfaces.CreateInvoiceTransactionComponent
	createNotificationComponent       notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction                mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                      string
}

func NewCreateInvoiceRepository(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	createInvoiceRepositoryTransactionComponent invoicedomainrepositoryinterfaces.CreateInvoiceTransactionComponent,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (invoicedomainrepositoryinterfaces.CreateInvoiceRepository, error) {
	createInvoiceRepo := &createInvoiceRepository{
		purchaseOrderDataSource,
		memberAccessDataSource,
		createInvoiceRepositoryTransactionComponent,
		createNotificationComponent,
		mongoDBTransaction,
		"CreateInvoiceRepository",
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
		input.(*invoicedomainrepositorytypes.CreateInvoiceInput),
	)
}

func (createInvoiceRepo *createInvoiceRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	invoiceToCreate := input.(*invoicedomainrepositorytypes.CreateInvoiceInput)

	return createInvoiceRepo.createInvoiceTransactionComponent.TransactionBody(
		operationOption,
		invoiceToCreate,
	)
}

func (createInvoiceRepo *createInvoiceRepository) RunTransaction(
	input *model.InternalCreateInvoice,
) ([]*model.Invoice, error) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	currentTime := time.Now().In(loc)
	if input.PaymentDueDate == nil {
		futureDateOnly := time.Date(
			currentTime.Year(),
			currentTime.Month(),
			currentTime.Day()+7,
			0, 0, 0, 0,
			currentTime.Location(),
		).UTC()
		input.PaymentDueDate = &futureDateOnly
	} else {
		dateOnly := time.Date(
			input.PaymentDueDate.Year(),
			input.PaymentDueDate.Month(),
			input.PaymentDueDate.Day(),
			0, 0, 0, 0,
			input.PaymentDueDate.Location(),
		).UTC()
		input.PaymentDueDate = &dateOnly
	}

	query := map[string]interface{}{
		"status": model.PurchaseOrderStatusWaitingForInvoice,
		"paymentDueDate": map[string]interface{}{
			"$lte": input.PaymentDueDate,
		},
	}
	if input.StartInvoiceDate != nil && input.EndInvoiceDate != nil {
		input.StartInvoiceDate = func(t time.Time) *time.Time { return &t }(
			time.Date(
				input.StartInvoiceDate.Year(),
				input.StartInvoiceDate.Month(),
				input.StartInvoiceDate.Day(),
				0, 0, 0, 0,
				input.StartInvoiceDate.Location(),
			).UTC(),
		)
		input.EndInvoiceDate = func(t time.Time) *time.Time { return &t }(
			time.Date(
				input.EndInvoiceDate.Year(),
				input.EndInvoiceDate.Month(),
				input.EndInvoiceDate.Day(),
				0, 0, 0, 0,
				input.EndInvoiceDate.Location(),
			).UTC(),
		)
		delete(query, "paymentDueDate")
		query["$and"] = []map[string]interface{}{
			{
				"paymentDueDate": map[string]interface{}{
					"$gte": input.StartInvoiceDate,
				},
			},
			{
				"paymentDueDate": map[string]interface{}{
					"$lte": input.EndInvoiceDate,
				},
			},
		}
		input.PaymentDueDate = input.EndInvoiceDate
	}

	purchaseOrders, err := createInvoiceRepo.purchaseOrderDataSource.GetMongoDataSource().Find(
		query,
		&mongodbcoretypes.PaginationOptions{},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createInvoiceRepo.pathIdentity,
			err,
		)
	}

	groupedPurchaseOrderByOrganization := map[string]map[string][]*model.PurchaseOrder{}
	for _, po := range purchaseOrders {
		orgStringID := po.Organization.ID.Hex()

		mouStringID := "NONE"
		if po.Mou != nil {
			mouStringID = po.Mou.ID.Hex()
		}

		if groupedPurchaseOrderByOrganization[orgStringID][mouStringID] == nil {
			groupedPurchaseOrderByOrganization[orgStringID][mouStringID] = []*model.PurchaseOrder{}
		}
		groupedPurchaseOrderByOrganization[orgStringID][mouStringID] = append(
			groupedPurchaseOrderByOrganization[orgStringID][mouStringID],
			po,
		)
	}

	createdInvoices := []*model.Invoice{}
	for _, orgKey := range reflect.ValueOf(groupedPurchaseOrderByOrganization).MapKeys() {
		for _, mouKey := range reflect.ValueOf(groupedPurchaseOrderByOrganization[orgKey.String()]).MapKeys() {
			purchaseOrders := groupedPurchaseOrderByOrganization[orgKey.String()][mouKey.String()]
			output, err := createInvoiceRepo.mongoDBTransaction.RunTransaction(
				&invoicedomainrepositorytypes.CreateInvoiceInput{
					CreateInvoiceInput:      input,
					PurchaseOrdersToInvoice: purchaseOrders,
				},
			)
			if err != nil {
				return nil, err
			}

			invoice := output.(*model.Invoice)
			createdInvoices = append(createdInvoices, invoice)
			go func() {
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
							InvoicePayload: &model.InvoicePayloadInput{
								Invoice: &model.InvoiceForNotifPayloadInput{},
							},
						},
					}
					json.Unmarshal(jsonInvPayload, &notifToCreate.PayloadOptions.InvoicePayload.Invoice)
					_, err := createInvoiceRepo.createNotificationComponent.TransactionBody(
						&mongodbcoretypes.OperationOptions{},
						notifToCreate,
					)
					if err != nil {
						return
					}
				}
			}()
		}
	}
	return createdInvoices, nil
}
