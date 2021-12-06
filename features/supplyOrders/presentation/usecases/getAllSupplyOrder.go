package supplyorderpresentationusecaseinterfaces

import (
	supplyorderpresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllSupplyOrderUsecase interface {
	Execute(input supplyorderpresentationusecasetypes.GetAllSupplyOrderUsecaseInput) ([]*model.SupplyOrder, error)
}
