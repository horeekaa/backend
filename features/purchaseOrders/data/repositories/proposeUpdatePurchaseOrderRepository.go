package purchaseorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	databasePurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseOrderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdatePurchaseOrderRepository struct {
	purchaseOrderDataSource                        databasePurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	proposeUpdatePurchaseOrderTransactionComponent purchaseOrderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderTransactionComponent
	createPurchaseOrderItemComponent               purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent
	updatePurchaseOrderItemComponent               purchaseorderitemdomainrepositoryinterfaces.UpdatePurchaseOrderItemTransactionComponent
	mongoDBTransaction                             mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdatePurchaseOrderRepository(
	purchaseOrderDataSource databasePurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	proposeUpdatePurchaseOrderRepositoryTransactionComponent purchaseOrderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderTransactionComponent,
	createPurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent,
	updatePurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.UpdatePurchaseOrderItemTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseOrderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderRepository, error) {
	proposeUpdatePurchaseOrderRepo := &proposeUpdatePurchaseOrderRepository{
		purchaseOrderDataSource,
		proposeUpdatePurchaseOrderRepositoryTransactionComponent,
		createPurchaseOrderItemComponent,
		updatePurchaseOrderItemComponent,
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

			purchaseOrderItemToCreate := &model.InternalCreatePurchaseOrderItem{}
			jsonTemp, _ := json.Marshal(purchaseOrderItemToUpdate)
			json.Unmarshal(jsonTemp, purchaseOrderItemToCreate)
			purchaseOrderItemToCreate.PurchaseOrder = &model.ObjectIDOnly{
				ID: &existingPurchaseOrder.ID,
			}

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
		if len(savedPurchaseOrderItems) > len(existingPurchaseOrder.Items) {
			jsonTemp, _ := json.Marshal(
				map[string]interface{}{
					"Items": savedPurchaseOrderItems,
				},
			)
			json.Unmarshal(jsonTemp, purchaseOrderToUpdate)
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
