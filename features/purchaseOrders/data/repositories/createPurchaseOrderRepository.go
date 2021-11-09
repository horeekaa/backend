package purchaseorderdomainrepositories

import (
	"encoding/json"
	"reflect"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createPurchaseOrderRepository struct {
	createPurchaseOrderTransactionComponent purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderTransactionComponent
	createPurchaseOrderItemComponent        purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent
	mongoDBTransaction                      mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreatePurchaseOrderRepository(
	createPurchaseOrderRepositoryTransactionComponent purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderTransactionComponent,
	createPurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderRepository, error) {
	createPurchaseOrderRepo := &createPurchaseOrderRepository{
		createPurchaseOrderRepositoryTransactionComponent,
		createPurchaseOrderItemComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createPurchaseOrderRepo,
		"CreatePurchaseOrderRepository",
	)

	return createPurchaseOrderRepo, nil
}

func (createPurchaseOrderRepo *createPurchaseOrderRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createPurchaseOrderRepo.createPurchaseOrderTransactionComponent.PreTransaction(
		input.(*model.InternalCreatePurchaseOrder),
	)
}

func (createPurchaseOrderRepo *createPurchaseOrderRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	purchaseOrderToCreate := input.(*model.InternalCreatePurchaseOrder)
	purchaseOrders := []*model.PurchaseOrder{}

	if len(purchaseOrderToCreate.MouItems) > 0 {
		generatedObjectID := createPurchaseOrderRepo.createPurchaseOrderTransactionComponent.GenerateNewObjectID()
		savedPurchaseOrderItems := map[string][]*model.InternalCreatePurchaseOrderItem{}
		mouForPurchaseOrders := map[string]*model.ObjectIDOnly{}
		for _, purchaseOrderItem := range purchaseOrderToCreate.MouItems {
			if purchaseOrderItem.MouItem == nil {
				continue
			}
			purchaseOrderItem.PurchaseOrder = &model.ObjectIDOnly{
				ID: &generatedObjectID,
			}
			createdPurchaseOrderOutput, err := createPurchaseOrderRepo.createPurchaseOrderItemComponent.TransactionBody(
				operationOption,
				purchaseOrderItem,
			)
			if err != nil {
				return nil, err
			}
			savedPurchaseOrderItem := &model.InternalCreatePurchaseOrderItem{}
			jsonTemp, _ := json.Marshal(createdPurchaseOrderOutput)
			json.Unmarshal(jsonTemp, savedPurchaseOrderItem)

			mouStringID := purchaseOrderItem.MouItem.Mou.ID.Hex()
			mouForPurchaseOrders[mouStringID] = purchaseOrderItem.MouItem.Mou

			if savedPurchaseOrderItems[mouStringID] == nil {
				savedPurchaseOrderItems[mouStringID] = []*model.InternalCreatePurchaseOrderItem{}
			}
			savedPurchaseOrderItems[mouStringID] = append(
				savedPurchaseOrderItems[mouStringID],
				savedPurchaseOrderItem,
			)
		}

		for _, e := range reflect.ValueOf(savedPurchaseOrderItems).MapKeys() {
			purchaseOrderToCreate.Items = savedPurchaseOrderItems[e.String()]
			purchaseOrderToCreate.Type = model.PurchaseOrderTypeMouBased
			purchaseOrderToCreate.Mou = &model.MouForPurchaseOrderInput{
				ID: *mouForPurchaseOrders[e.String()].ID,
			}

			purchaseOrder, err := createPurchaseOrderRepo.createPurchaseOrderTransactionComponent.TransactionBody(
				operationOption,
				purchaseOrderToCreate,
			)
			if err != nil {
				return nil, err
			}

			purchaseOrders = append(purchaseOrders, purchaseOrder)
		}
	}

	if len(purchaseOrderToCreate.RetailItems) > 0 {
		generatedObjectID := createPurchaseOrderRepo.createPurchaseOrderTransactionComponent.GenerateNewObjectID()
		savedPurchaseOrderItems := []*model.InternalCreatePurchaseOrderItem{}
		for _, purchaseOrderItem := range purchaseOrderToCreate.RetailItems {
			purchaseOrderItem.PurchaseOrder = &model.ObjectIDOnly{
				ID: &generatedObjectID,
			}
			createdPurchaseOrderOutput, err := createPurchaseOrderRepo.createPurchaseOrderItemComponent.TransactionBody(
				operationOption,
				purchaseOrderItem,
			)
			if err != nil {
				return nil, err
			}
			savedPurchaseOrderItem := &model.InternalCreatePurchaseOrderItem{}
			jsonTemp, _ := json.Marshal(createdPurchaseOrderOutput)
			json.Unmarshal(jsonTemp, savedPurchaseOrderItem)
			savedPurchaseOrderItems = append(savedPurchaseOrderItems, savedPurchaseOrderItem)
		}
		purchaseOrderToCreate.Items = savedPurchaseOrderItems
		purchaseOrderToCreate.Type = model.PurchaseOrderTypeRetail
		purchaseOrder, err := createPurchaseOrderRepo.createPurchaseOrderTransactionComponent.TransactionBody(
			operationOption,
			purchaseOrderToCreate,
		)
		if err != nil {
			return nil, err
		}

		purchaseOrders = append(purchaseOrders, purchaseOrder)
	}

	return purchaseOrders, nil
}

func (createPurchaseOrderRepo *createPurchaseOrderRepository) RunTransaction(
	input *model.InternalCreatePurchaseOrder,
) ([]*model.PurchaseOrder, error) {
	output, err := createPurchaseOrderRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).([]*model.PurchaseOrder), nil
}
