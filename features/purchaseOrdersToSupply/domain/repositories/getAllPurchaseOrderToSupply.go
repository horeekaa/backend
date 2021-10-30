package purchaseordertosupplydomainrepositoryinterfaces

import (
	purchaseordertosupplydomainrepositorytypes "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllPurchaseOrderToSupplyRepository interface {
	Execute(filterFields purchaseordertosupplydomainrepositorytypes.GetAllPurchaseOrderToSupplyInput) ([]*model.PurchaseOrderToSupply, error)
}
