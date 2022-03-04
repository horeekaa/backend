package supplyorderitempresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type UpdateSupplyOrderItemPickUpUsecaseInput struct {
	Context                     context.Context
	UpdateSupplyOrderItemPickUp *model.UpdateSupplyOrderItemPickUpOnly
}

type GetAllSupplyOrderItemUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.SupplyOrderItemFilterFields
	PaginationOps *model.PaginationOptionInput
}
