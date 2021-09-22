package purchaseorderitemdomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type PurchaseOrderItemLoader interface {
	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		mouItem *model.MouItemForPurchaseOrderItemInput,
		productVariant *model.ProductVariantForPurchaseOrderItemInput,
	) (bool, error)
}
