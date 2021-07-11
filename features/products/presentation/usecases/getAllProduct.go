package productpresentationusecaseinterfaces

import (
	productpresentationusecasetypes "github.com/horeekaa/backend/features/products/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllProductUsecase interface {
	Execute(input productpresentationusecasetypes.GetAllProductUsecaseInput) ([]*model.Product, error)
}
