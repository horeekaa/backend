package purchaseorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdatePurchaseOrderRepository struct {
	purchaseOrderDataSource                        databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	proposeUpdatePurchaseOrderTransactionComponent purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderTransactionComponent
	createPurchaseOrderItemComponent               purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent
	updatePurchaseOrderItemComponent               purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent
	approvePurchaseOrderItemComponent              purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent
	updateInvoiceTrxComponent                      invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent
	mongoDBTransaction                             mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                                   string
}

func NewProposeUpdatePurchaseOrderRepository(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	proposeUpdatePurchaseOrderRepositoryTransactionComponent purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderTransactionComponent,
	createPurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent,
	updatePurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent,
	approvePurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent,
	updateInvoiceTrxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderRepository, error) {
	proposeUpdatePurchaseOrderRepo := &proposeUpdatePurchaseOrderRepository{
		purchaseOrderDataSource,
		proposeUpdatePurchaseOrderRepositoryTransactionComponent,
		createPurchaseOrderItemComponent,
		updatePurchaseOrderItemComponent,
		approvePurchaseOrderItemComponent,
		updateInvoiceTrxComponent,
		mongoDBTransaction,
		"ProposeUpdatePurchaseOrderRepository",
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdatePurchaseOrderRepo,
		"ProposeUpdatePurchaseOrderRepository",
	)

	return proposeUpdatePurchaseOrderRepo, nil
}

func (updatePurchaseOrderRepo *proposeUpdatePurchaseOrderRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updatePurchaseOrderRepo.proposeUpdatePurchaseOrderTransactionComponent.PreTransaction(
		input.(*model.InternalUpdatePurchaseOrder),
	)
}

func (updatePurchaseOrderRepo *proposeUpdatePurchaseOrderRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	purchaseOrderToUpdate := input.(*model.InternalUpdatePurchaseOrder)
	existingPurchaseOrder, err := updatePurchaseOrderRepo.purchaseOrderDataSource.GetMongoDataSource().FindByID(
		purchaseOrderToUpdate.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updatePurchaseOrderRepo.pathIdentity,
			err,
		)
	}

	if purchaseOrderToUpdate.Items != nil {
		savedPurchaseOrderItems := existingPurchaseOrder.Items
		for _, purchaseOrderItemToUpdate := range purchaseOrderToUpdate.Items {
			if purchaseOrderItemToUpdate.ID != nil {
				if !funk.Contains(
					existingPurchaseOrder.Items,
					func(mi *model.PurchaseOrderItem) bool {
						return mi.ID == *purchaseOrderItemToUpdate.ID
					},
				) {
					continue
				}
				if purchaseOrderItemToUpdate.ProposalStatus != nil {
					purchaseOrderItemToUpdate.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
						return &m
					}(*purchaseOrderToUpdate.SubmittingAccount)

					_, err := updatePurchaseOrderRepo.approvePurchaseOrderItemComponent.TransactionBody(
						operationOption,
						purchaseOrderItemToUpdate,
					)
					if err != nil {
						return nil, err
					}
					continue
				}

				purchaseOrderItemToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*purchaseOrderToUpdate.ProposalStatus)
				purchaseOrderItemToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*purchaseOrderToUpdate.SubmittingAccount)
				if *purchaseOrderToUpdate.MemberAccess.Organization.Type != model.OrganizationTypeCustomer {
					purchaseOrderItemToUpdate.CustomerAgreed = func(b bool) *bool { return &b }(false)
				}

				_, err := updatePurchaseOrderRepo.updatePurchaseOrderItemComponent.TransactionBody(
					operationOption,
					purchaseOrderItemToUpdate,
				)
				if err != nil {
					return nil, err
				}
				continue
			}
			if existingPurchaseOrder.Mou != nil && purchaseOrderItemToUpdate.MouItem == nil {
				return nil, horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.POItemMismatchWithPOType,
					updatePurchaseOrderRepo.pathIdentity,
					nil,
				)
			}

			purchaseOrderItemToCreate := &model.InternalCreatePurchaseOrderItem{}
			jsonTemp, _ := json.Marshal(purchaseOrderItemToUpdate)
			json.Unmarshal(jsonTemp, purchaseOrderItemToCreate)
			purchaseOrderItemToCreate.PurchaseOrder = &model.ObjectIDOnly{
				ID: &existingPurchaseOrder.ID,
			}
			purchaseOrderItemToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*purchaseOrderToUpdate.ProposalStatus)
			purchaseOrderItemToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*purchaseOrderToUpdate.SubmittingAccount)
			if *purchaseOrderToUpdate.MemberAccess.Organization.Type == model.OrganizationTypeCustomer {
				purchaseOrderItemToCreate.CustomerAgreed = func(b bool) *bool { return &b }(true)
			}

			savedPurchaseOrderItem, err := updatePurchaseOrderRepo.createPurchaseOrderItemComponent.TransactionBody(
				operationOption,
				purchaseOrderItemToCreate,
			)
			if err != nil {
				return nil, err
			}
			savedPurchaseOrderItems = append(savedPurchaseOrderItems, savedPurchaseOrderItem)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"Items": savedPurchaseOrderItems,
			},
		)
		json.Unmarshal(jsonTemp, purchaseOrderToUpdate)
	}

	if purchaseOrderToUpdate.ProposalStatus != nil {
		if *purchaseOrderToUpdate.ProposalStatus == model.EntityProposalStatusApproved {
			if purchaseOrderToUpdate.Invoice != nil {
				if funk.Get(existingPurchaseOrder, "Invoice.ID") != nil {
					if existingPurchaseOrder.Invoice.ID.Hex() != purchaseOrderToUpdate.Invoice.ID.Hex() {
						_, err := updatePurchaseOrderRepo.updateInvoiceTrxComponent.TransactionBody(
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
				_, err := updatePurchaseOrderRepo.updateInvoiceTrxComponent.TransactionBody(
					operationOption,
					&model.InternalUpdateInvoice{
						ID: *purchaseOrderToUpdate.Invoice.ID,
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

	return updatePurchaseOrderRepo.proposeUpdatePurchaseOrderTransactionComponent.TransactionBody(
		operationOption,
		purchaseOrderToUpdate,
	)
}

func (updatePurchaseOrderRepo *proposeUpdatePurchaseOrderRepository) RunTransaction(
	input *model.InternalUpdatePurchaseOrder,
) (*model.PurchaseOrder, error) {
	output, err := updatePurchaseOrderRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.PurchaseOrder), nil
}
