package supplyorderitempresentationusecaseinterfaces

import (
	supplyorderitempresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllSupplyOrderItemUsecase interface {
	Execute(input supplyorderitempresentationusecasetypes.GetAllSupplyOrderItemUsecaseInput) ([]*model.SupplyOrderItem, error)
}
