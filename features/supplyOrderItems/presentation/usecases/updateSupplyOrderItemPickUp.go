package supplyorderitempresentationusecaseinterfaces

import (
	supplyorderitempresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdateSupplyOrderItemPickUpUsecase interface {
	Execute(input supplyorderitempresentationusecasetypes.UpdateSupplyOrderItemPickUpUsecaseInput) (*model.SupplyOrderItem, error)
}
