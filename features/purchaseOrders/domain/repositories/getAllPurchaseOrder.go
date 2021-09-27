package purchaseorderdomainrepositoryinterfaces

import (
	purchaseorderdomainrepositorytypes "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllPurchaseOrderRepository interface {
	Execute(filterFields purchaseorderdomainrepositorytypes.GetAllPurchaseOrderInput) ([]*model.PurchaseOrder, error)
}
