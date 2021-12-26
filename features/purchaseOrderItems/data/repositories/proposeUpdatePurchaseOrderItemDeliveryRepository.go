package purchaseorderitemdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdatePurchaseOrderItemDeliveryRepository struct {
	purchaseOrderItemDataSource                        databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	proposeUpdatePurchaseOrderItemTransactionComponent purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent
	mongoDBTransaction                                 mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdatePurchaseOrderItemDeliveryRepository(
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
	proposeUpdatePurchaseOrderItemTransactionComponent purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemDeliveryRepository, error) {
	proposeUpdatePurchaseOrderItemRepo := &proposeUpdatePurchaseOrderItemDeliveryRepository{
		purchaseOrderItemDataSource,
		proposeUpdatePurchaseOrderItemTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdatePurchaseOrderItemRepo,
		"ProposeUpdatePurchaseOrderItemDeliveryRepository",
	)

	return proposeUpdatePurchaseOrderItemRepo, nil
}

func (updatePurchaseOrderItemRepo *proposeUpdatePurchaseOrderItemDeliveryRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updatePurchaseOrderItemRepo.proposeUpdatePurchaseOrderItemTransactionComponent.PreTransaction(
		input.(*model.InternalUpdatePurchaseOrderItem),
	)
}

func (updatePurchaseOrderItemRepo *proposeUpdatePurchaseOrderItemDeliveryRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	poItemToUpdate := input.(*model.InternalUpdatePurchaseOrderItem)

	return updatePurchaseOrderItemRepo.proposeUpdatePurchaseOrderItemTransactionComponent.TransactionBody(
		operationOption,
		poItemToUpdate,
	)
}

func (updatePurchaseOrderItemRepo *proposeUpdatePurchaseOrderItemDeliveryRepository) RunTransaction(
	input *model.InternalUpdatePurchaseOrderItem,
) (*model.PurchaseOrderItem, error) {
	output, err := updatePurchaseOrderItemRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	purchaseOrderItem := (output).(*model.PurchaseOrderItem)
	return purchaseOrderItem, nil
}
