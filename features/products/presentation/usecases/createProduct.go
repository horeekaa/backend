package productpresentationusecaseinterfaces

import (
	productpresentationusecasetypes "github.com/horeekaa/backend/features/products/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreateProductUsecase interface {
	Execute(input productpresentationusecasetypes.CreateProductUsecaseInput) (*model.Product, error)
}
