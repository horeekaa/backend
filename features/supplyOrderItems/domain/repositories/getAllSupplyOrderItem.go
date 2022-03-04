package supplyorderitemdomainrepositoryinterfaces

import (
	supplyorderitemdomainrepositorytypes "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllSupplyOrderItemRepository interface {
	Execute(filterFields supplyorderitemdomainrepositorytypes.GetAllSupplyOrderItemInput) ([]*model.SupplyOrderItem, error)
}
