package purchaseorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
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
	mongoDBTransaction                             mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdatePurchaseOrderRepository(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	proposeUpdatePurchaseOrderRepositoryTransactionComponent purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderTransactionComponent,
	createPurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent,
	updatePurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent,
	approvePurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderRepository, error) {
	proposeUpdatePurchaseOrderRepo := &proposeUpdatePurchaseOrderRepository{
		purchaseOrderDataSource,
		proposeUpdatePurchaseOrderRepositoryTransactionComponent,
		createPurchaseOrderItemComponent,
		updatePurchaseOrderItemComponent,
		approvePurchaseOrderItemComponent,
		mongoDBTransaction,
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
			"/proposeUpdatePurchaseOrderRepository",
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
						return nil, horeekaacoreexceptiontofailure.ConvertException(
							"/proposeUpdatepurchaseOrderRepository",
							err,
						)
					}
					continue
				}

				purchaseOrderItemToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*purchaseOrderToUpdate.ProposalStatus)
				purchaseOrderItemToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*purchaseOrderToUpdate.SubmittingAccount)

				_, err := updatePurchaseOrderRepo.updatePurchaseOrderItemComponent.TransactionBody(
					operationOption,
					purchaseOrderItemToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdatePurchaseOrderRepository",
						err,
					)
				}
				continue
			}
			if existingPurchaseOrder.Mou != nil && purchaseOrderItemToUpdate.MouItem == nil {
				return nil, horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.POItemMismatchWithPOType,
					"/proposeUpdatePurchaseOrderRepository",
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

			savedPurchaseOrderItem, err := updatePurchaseOrderRepo.createPurchaseOrderItemComponent.TransactionBody(
				operationOption,
				purchaseOrderItemToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/proposeUpdatePurchaseOrderRepository",
					err,
				)
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
