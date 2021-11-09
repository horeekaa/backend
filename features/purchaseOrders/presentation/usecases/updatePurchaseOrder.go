package purchaseorderpresentationusecaseinterfaces

import (
	purchaseorderpresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdatePurchaseOrderUsecase interface {
	Execute(input purchaseorderpresentationusecasetypes.UpdatePurchaseOrderUsecaseInput) (*model.PurchaseOrder, error)
}
