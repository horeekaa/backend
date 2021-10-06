package purchaseorderitempresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetPurchaseOrderItemUsecase interface {
	Execute(input *model.PurchaseOrderItemFilterFields) (*model.PurchaseOrderItem, error)
}
