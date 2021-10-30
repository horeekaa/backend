package purchaseordertosupplypresentationusecaseinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type CreatePurchaseOrderToSupplyUsecase interface {
	Execute() ([]*model.PurchaseOrderToSupply, error)
}
