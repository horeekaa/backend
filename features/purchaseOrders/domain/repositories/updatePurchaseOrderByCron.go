package purchaseorderdomainrepositoryinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type UpdatePurchaseOrderByCronRepository interface {
	RunTransaction() ([]*model.PurchaseOrder, error)
}
