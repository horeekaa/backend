package purchaseordertosupplydomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type CreatePurchaseOrderToSupplyTransactionComponent interface {
	PreTransaction(
		createPurchaseOrderToSupplyInput *model.PurchaseOrder,
	) (*model.PurchaseOrder, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createPurchaseOrderToSupplyInput *model.PurchaseOrder,
	) ([]*model.PurchaseOrderToSupply, error)
}

type CreatePurchaseOrderToSupplyRepository interface {
	RunTransaction() ([]*model.PurchaseOrderToSupply, error)
}
