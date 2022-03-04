package purchaseorderitempresentationusecaseinterfaces

import (
	purchaseorderitempresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllPurchaseOrderItemUsecase interface {
	Execute(input purchaseorderitempresentationusecasetypes.GetAllPurchaseOrderItemUsecaseInput) ([]*model.PurchaseOrderItem, error)
}
