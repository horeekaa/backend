package supplyorderpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetSupplyOrderUsecase interface {
	Execute(input *model.SupplyOrderFilterFields) (*model.SupplyOrder, error)
}
