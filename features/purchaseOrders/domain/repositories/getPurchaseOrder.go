package purchaseorderdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetPurchaseOrderRepository interface {
	Execute(filterFields *model.PurchaseOrderFilterFields) (*model.PurchaseOrder, error)
}
