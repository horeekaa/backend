package productpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreateProductUsecaseInput struct {
	Context       context.Context
	CreateProduct *model.CreateProduct
}

type UpdateProductUsecaseInput struct {
	Context       context.Context
	UpdateProduct *model.UpdateProduct
}

type GetAllProductUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.ProductFilterFields
	PaginationOps *model.PaginationOptionInput
}
