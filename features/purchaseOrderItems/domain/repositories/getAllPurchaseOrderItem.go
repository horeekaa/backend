package purchaseorderitemdomainrepositoryinterfaces

import (
	purchaseorderitemdomainrepositorytypes "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllPurchaseOrderItemRepository interface {
	Execute(filterFields purchaseorderitemdomainrepositorytypes.GetAllPurchaseOrderItemInput) ([]*model.PurchaseOrderItem, error)
}
