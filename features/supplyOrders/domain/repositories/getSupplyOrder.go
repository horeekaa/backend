package supplyorderdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetSupplyOrderRepository interface {
	Execute(filterFields *model.SupplyOrderFilterFields) (*model.SupplyOrder, error)
}
