package purchaseordertosupplydomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetPurchaseOrderToSupplyRepository interface {
	Execute(filterFields *model.PurchaseOrderToSupplyFilterFields) (*model.PurchaseOrderToSupply, error)
}
