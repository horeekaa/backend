package supplyorderpresentationusecaseinterfaces

import (
	supplyorderpresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreateSupplyOrderUsecase interface {
	Execute(input supplyorderpresentationusecasetypes.CreateSupplyOrderUsecaseInput) ([]*model.SupplyOrder, error)
}
