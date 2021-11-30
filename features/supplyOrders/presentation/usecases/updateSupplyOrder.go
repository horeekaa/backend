package supplyorderpresentationusecaseinterfaces

import (
	supplyorderpresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdateSupplyOrderUsecase interface {
	Execute(input supplyorderpresentationusecasetypes.UpdateSupplyOrderUsecaseInput) (*model.SupplyOrder, error)
}
