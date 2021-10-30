package purchaseordertosupplypresentationusecaseinterfaces

import (
	purchaseordertosupplypresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllPurchaseOrderToSupplyUsecase interface {
	Execute(input purchaseordertosupplypresentationusecasetypes.GetAllPurchaseOrderToSupplyUsecaseInput) ([]*model.PurchaseOrderToSupply, error)
}
