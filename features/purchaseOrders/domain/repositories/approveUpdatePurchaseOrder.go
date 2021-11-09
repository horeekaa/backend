package purchaseorderdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ApproveUpdatePurchaseOrderTransactionComponent interface {
	PreTransaction(
		updatePurchaseOrderInput *model.InternalUpdatePurchaseOrder,
	) (*model.InternalUpdatePurchaseOrder, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updatePurchaseOrderInput *model.InternalUpdatePurchaseOrder,
	) (*model.PurchaseOrder, error)
}

type ApproveUpdatePurchaseOrderRepository interface {
	RunTransaction(
		updatePurchaseOrderInput *model.InternalUpdatePurchaseOrder,
	) (*model.PurchaseOrder, error)
}
