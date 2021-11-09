package purchaseordertosupplydomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ProcessPurchaseOrderToSupplyTransactionComponent interface {
	PreTransaction(
		input *model.PurchaseOrderToSupply,
	) (*model.PurchaseOrderToSupply, error)

	TransactionBody(
		operationOption *mongodbcoretypes.OperationOptions,
		input *model.PurchaseOrderToSupply,
	) ([]*model.InternalCreateNotification, error)
}

type ProcessPurchaseOrderToSupplyRepository interface {
	RunTransaction() (bool, error)
}
