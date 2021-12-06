package supplyorderpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreateSupplyOrderUsecaseInput struct {
	Context           context.Context
	CreateSupplyOrder *model.CreateSupplyOrder
}

type UpdateSupplyOrderUsecaseInput struct {
	Context           context.Context
	UpdateSupplyOrder *model.UpdateSupplyOrder
}

type GetAllSupplyOrderUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.SupplyOrderFilterFields
	PaginationOps *model.PaginationOptionInput
}
