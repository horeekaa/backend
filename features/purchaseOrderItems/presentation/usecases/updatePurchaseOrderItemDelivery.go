package purchaseorderitempresentationusecaseinterfaces

import (
	purchaseorderitempresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdatePurchaseOrderItemDeliveryUsecase interface {
	Execute(input purchaseorderitempresentationusecasetypes.UpdatePurchaseOrderItemDeliveryUsecaseInput) (*model.PurchaseOrderItem, error)
}
