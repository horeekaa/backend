package supplyorderdomainrepositoryinterfaces

import (
	supplyorderdomainrepositorytypes "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllSupplyOrderRepository interface {
	Execute(filterFields supplyorderdomainrepositorytypes.GetAllSupplyOrderInput) ([]*model.SupplyOrder, error)
}
