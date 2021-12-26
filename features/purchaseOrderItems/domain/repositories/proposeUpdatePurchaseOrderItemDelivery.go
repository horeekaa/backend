package purchaseorderitemdomainrepositoryinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type ProposeUpdatePurchaseOrderItemDeliveryRepository interface {
	RunTransaction(
		updatePurchaseOrderItemInput *model.InternalUpdatePurchaseOrderItem,
	) (*model.PurchaseOrderItem, error)
}
