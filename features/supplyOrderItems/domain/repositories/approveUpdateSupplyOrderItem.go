package supplyorderitemdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ApproveUpdateSupplyOrderItemTransactionComponent interface {
	PreTransaction(
		updateSupplyOrderItemInput *model.InternalUpdateSupplyOrderItem,
	) (*model.InternalUpdateSupplyOrderItem, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateSupplyOrderItemInput *model.InternalUpdateSupplyOrderItem,
	) (*model.SupplyOrderItem, error)
}
