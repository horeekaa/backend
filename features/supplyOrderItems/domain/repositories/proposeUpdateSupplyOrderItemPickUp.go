package supplyorderitemdomainrepositoryinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type ProposeUpdateSupplyOrderItemPickUpRepository interface {
	RunTransaction(
		updateSupplyOrderItemInput *model.InternalUpdateSupplyOrderItem,
	) (*model.SupplyOrderItem, error)
}
