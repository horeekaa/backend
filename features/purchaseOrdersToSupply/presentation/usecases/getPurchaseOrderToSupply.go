package purchaseordertosupplypresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetPurchaseOrderToSupplyUsecase interface {
	Execute(input *model.PurchaseOrderToSupplyFilterFields) (*model.PurchaseOrderToSupply, error)
}
