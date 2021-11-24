package supplyorderdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ApproveUpdateSupplyOrderTransactionComponent interface {
	PreTransaction(
		updateSupplyOrderInput *model.InternalUpdateSupplyOrder,
	) (*model.InternalUpdateSupplyOrder, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateSupplyOrderInput *model.InternalUpdateSupplyOrder,
	) (*model.SupplyOrder, error)
}

type ApproveUpdateSupplyOrderRepository interface {
	RunTransaction(
		updateSupplyOrderInput *model.InternalUpdateSupplyOrder,
	) (*model.SupplyOrder, error)
}
