package purchaseorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type approveUpdatePurchaseOrderRepository struct {
	memberAccessDataSource                         databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	purchaseOrderDataSource                        databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	approveUpdatePurchaseOrderItemComponent        purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent
	approveUpdatePurchaseOrderTransactionComponent purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderTransactionComponent
	updateInvoiceTrxComponent                      invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent
	createNotificationComponent                    notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction                             mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                                   string
}

func NewApproveUpdatePurchaseOrderRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	approveUpdatePurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent,
	approveUpdatepurchaseOrderTransactionComponent purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderTransactionComponent,
	updateInvoiceTrxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderRepository, error) {
	approveUpdatePurchaseOrderRepo := &approveUpdatePurchaseOrderRepository{
		memberAccessDataSource,
		purchaseOrderDataSource,
		approveUpdatePurchaseOrderItemComponent,
		approveUpdatepurchaseOrderTransactionComponent,
		updateInvoiceTrxComponent,
		createNotificationComponent,
		mongoDBTransaction,
		"ApproveUpdatePurchaseOrderRepository",
	}

	mongoDBTransaction.SetTransaction(
		approveUpdatePurchaseOrderRepo,
		"ApproveUpdatePurchaseOrderRepository",
	)

	return approveUpdatePurchaseOrderRepo, nil
}

func (approveUpdatePurchaseOrderRepo *approveUpdatePurchaseOrderRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return input, nil
}

func (approveUpdatePurchaseOrderRepo *approveUpdatePurchaseOrderRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	purchaseOrderToApprove := input.(*model.InternalUpdatePurchaseOrder)
	existingPurchaseOrder, err := approveUpdatePurchaseOrderRepo.purchaseOrderDataSource.GetMongoDataSource().FindByID(
		purchaseOrderToApprove.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveUpdatePurchaseOrderRepo.pathIdentity,
			err,
		)
	}
	if existingPurchaseOrder.ProposedChanges.ProposalStatus == model.EntityProposalStatusProposed {
		if existingPurchaseOrder.ProposedChanges.Items != nil {
			for _, poItem := range existingPurchaseOrder.ProposedChanges.Items {
				updatePurchaseOrderItem := &model.InternalUpdatePurchaseOrderItem{
					ID: &poItem.ID,
				}
				updatePurchaseOrderItem.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*purchaseOrderToApprove.RecentApprovingAccount)
				updatePurchaseOrderItem.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*purchaseOrderToApprove.ProposalStatus)

				_, err := approveUpdatePurchaseOrderRepo.approveUpdatePurchaseOrderItemComponent.TransactionBody(
					operationOption,
					updatePurchaseOrderItem,
				)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if purchaseOrderToApprove.ProposalStatus != nil {
		if *purchaseOrderToApprove.ProposalStatus == model.EntityProposalStatusApproved {
			if existingPurchaseOrder.ProposedChanges.Invoice != nil {
				if funk.Get(existingPurchaseOrder, "Invoice.ID") != nil {
					if existingPurchaseOrder.Invoice.ID.Hex() != existingPurchaseOrder.ProposedChanges.Invoice.ID.Hex() {
						_, err := approveUpdatePurchaseOrderRepo.updateInvoiceTrxComponent.TransactionBody(
							operationOption,
							&model.InternalUpdateInvoice{
								ID: *existingPurchaseOrder.Invoice.ID,
								PurchaseOrdersToRemove: []*model.ObjectIDOnly{
									{ID: &existingPurchaseOrder.ID},
								},
							},
						)
						if err != nil {
							return nil, err
						}
					}
				}
				_, err := approveUpdatePurchaseOrderRepo.updateInvoiceTrxComponent.TransactionBody(
					operationOption,
					&model.InternalUpdateInvoice{
						ID: *existingPurchaseOrder.ProposedChanges.Invoice.ID,
						PurchaseOrdersToAdd: []*model.ObjectIDOnly{
							{ID: &existingPurchaseOrder.ID},
						},
					},
				)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return approveUpdatePurchaseOrderRepo.approveUpdatePurchaseOrderTransactionComponent.TransactionBody(
		operationOption,
		purchaseOrderToApprove,
	)
}

func (approveUpdatePurchaseOrderRepo *approveUpdatePurchaseOrderRepository) RunTransaction(
	input *model.InternalUpdatePurchaseOrder,
) (*model.PurchaseOrder, error) {
	output, err := approveUpdatePurchaseOrderRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	approvedPurchaseOrder := (output).(*model.PurchaseOrder)
	go func() {
		memberAccesses, err := approveUpdatePurchaseOrderRepo.memberAccessDataSource.GetMongoDataSource().Find(
			map[string]interface{}{
				"memberAccessRefType": model.MemberAccessRefTypeOrganizationsBased,
				"status":              model.MemberAccessStatusActive,
				"proposalStatus":      model.EntityProposalStatusApproved,
				"invitationAccepted":  true,
				"organization._id":    approvedPurchaseOrder.Organization.ID,
			},
			&mongodbcoretypes.PaginationOptions{},
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			return
		}

		for _, memberAccess := range memberAccesses {
			notificationToCreate := &model.InternalCreateNotification{
				NotificationCategory: model.NotificationCategoryPurchaseOrderApproval,
				PayloadOptions: &model.PayloadOptionsInput{
					PurchaseOrderPayload: &model.PurchaseOrderPayloadInput{
						PurchaseOrder: &model.PurchaseOrderForNotifPayloadInput{},
					},
				},
				RecipientAccount: &model.ObjectIDOnly{
					ID: &memberAccess.Account.ID,
				},
			}

			jsonTemp, _ := json.Marshal(approvedPurchaseOrder)
			json.Unmarshal(jsonTemp, &notificationToCreate.PayloadOptions.PurchaseOrderPayload.PurchaseOrder)

			_, err = approveUpdatePurchaseOrderRepo.createNotificationComponent.TransactionBody(
				&mongodbcoretypes.OperationOptions{},
				notificationToCreate,
			)
			if err != nil {
				return
			}
		}
	}()

	return approvedPurchaseOrder, err
}
