package productpresentationusecaseinterfaces

import (
	productpresentationusecasetypes "github.com/horeekaa/backend/features/products/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdateProductUsecase interface {
	Execute(input productpresentationusecasetypes.UpdateProductUsecaseInput) (*model.Product, error)
}
