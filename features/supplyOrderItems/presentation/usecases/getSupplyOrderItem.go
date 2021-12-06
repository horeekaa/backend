package supplyorderitempresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetSupplyOrderItemUsecase interface {
	Execute(input *model.SupplyOrderItemFilterFields) (*model.SupplyOrderItem, error)
}
