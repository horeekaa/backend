package supplyorderitemdomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type SupplyOrderItemLoader interface {
	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		purchaseOrderToSupply *model.PurchaseOrderToSupplyForSupplyOrderItemInput,
	) (bool, error)
}
