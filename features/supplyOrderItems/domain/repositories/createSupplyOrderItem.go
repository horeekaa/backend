package supplyorderitemdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type CreateSupplyOrderItemTransactionComponent interface {
	PreTransaction(
		createSupplyOrderItemInput *model.InternalCreateSupplyOrderItem,
	) (*model.InternalCreateSupplyOrderItem, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createSupplyOrderItemInput *model.InternalCreateSupplyOrderItem,
	) (*model.SupplyOrderItem, error)
}
