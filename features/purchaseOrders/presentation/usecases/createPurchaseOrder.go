package purchaseorderpresentationusecaseinterfaces

import (
	purchaseorderpresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreatePurchaseOrderUsecase interface {
	Execute(input purchaseorderpresentationusecasetypes.CreatePurchaseOrderUsecaseInput) ([]*model.PurchaseOrder, error)
}
