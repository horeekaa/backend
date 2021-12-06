package supplyorderitemdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetSupplyOrderItemRepository interface {
	Execute(filterFields *model.SupplyOrderItemFilterFields) (*model.SupplyOrderItem, error)
}
