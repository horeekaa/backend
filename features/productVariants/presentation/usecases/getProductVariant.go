package productvariantpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetProductVariantUsecase interface {
	Execute(input *model.ProductVariantFilterFields) (*model.ProductVariant, error)
}
