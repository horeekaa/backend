package purchaseorderitemdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type CreatePurchaseOrderItemTransactionComponent interface {
	PreTransaction(
		createPurchaseOrderItemInput *model.InternalCreatePurchaseOrderItem,
	) (*model.InternalCreatePurchaseOrderItem, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createPurchaseOrderItemInput *model.InternalCreatePurchaseOrderItem,
	) (*model.PurchaseOrderItem, error)
}
