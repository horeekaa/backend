package supplyorderitemdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateSupplyOrderItemPickUpRepository struct {
	supplyOrderItemDataSource                        databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource
	proposeUpdateSupplyOrderItemTransactionComponent supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemTransactionComponent
	mongoDBTransaction                               mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateSupplyOrderItemPickUpRepository(
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
	proposeUpdateSupplyOrderItemTransactionComponent supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemPickUpRepository, error) {
	proposeUpdateSupplyOrderItemRepo := &proposeUpdateSupplyOrderItemPickUpRepository{
		supplyOrderItemDataSource,
		proposeUpdateSupplyOrderItemTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateSupplyOrderItemRepo,
		"ProposeUpdateSupplyOrderItemPickUpRepository",
	)

	return proposeUpdateSupplyOrderItemRepo, nil
}

func (updateSupplyOrderItemRepo *proposeUpdateSupplyOrderItemPickUpRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateSupplyOrderItemRepo.proposeUpdateSupplyOrderItemTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateSupplyOrderItem),
	)
}

func (updateSupplyOrderItemRepo *proposeUpdateSupplyOrderItemPickUpRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	poItemToUpdate := input.(*model.InternalUpdateSupplyOrderItem)

	return updateSupplyOrderItemRepo.proposeUpdateSupplyOrderItemTransactionComponent.TransactionBody(
		operationOption,
		poItemToUpdate,
	)
}

func (updateSupplyOrderItemRepo *proposeUpdateSupplyOrderItemPickUpRepository) RunTransaction(
	input *model.InternalUpdateSupplyOrderItem,
) (*model.SupplyOrderItem, error) {
	output, err := updateSupplyOrderItemRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	supplyOrderItem := (output).(*model.SupplyOrderItem)
	return supplyOrderItem, nil
}
