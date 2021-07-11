package productpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetProductUsecase interface {
	Execute(input *model.ProductFilterFields) (*model.Product, error)
}
