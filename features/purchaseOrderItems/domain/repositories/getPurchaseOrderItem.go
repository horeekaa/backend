package purchaseorderitemdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetPurchaseOrderItemRepository interface {
	Execute(filterFields *model.PurchaseOrderItemFilterFields) (*model.PurchaseOrderItem, error)
}
