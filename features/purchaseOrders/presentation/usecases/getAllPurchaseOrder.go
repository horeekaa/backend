package purchaseorderpresentationusecaseinterfaces

import (
	purchaseOrderpresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllPurchaseOrderUsecase interface {
	Execute(input purchaseOrderpresentationusecasetypes.GetAllPurchaseOrderUsecaseInput) ([]*model.PurchaseOrder, error)
}
