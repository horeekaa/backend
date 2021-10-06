package purchaseorderpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetPurchaseOrderUsecase interface {
	Execute(input *model.PurchaseOrderFilterFields) (*model.PurchaseOrder, error)
}
