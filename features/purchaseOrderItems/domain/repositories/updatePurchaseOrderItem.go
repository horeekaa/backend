package purchaseorderitemdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type UpdatePurchaseOrderItemTransactionComponent interface {
	PreTransaction(
		updatePurchaseOrderItemInput *model.InternalUpdatePurchaseOrderItem,
	) (*model.InternalUpdatePurchaseOrderItem, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updatePurchaseOrderItemInput *model.InternalUpdatePurchaseOrderItem,
	) (*model.PurchaseOrderItem, error)
}
