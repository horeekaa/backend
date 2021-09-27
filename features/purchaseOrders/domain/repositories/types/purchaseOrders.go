package purchaseorderdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllPurchaseOrderInput struct {
	FilterFields  *model.PurchaseOrderFilterFields
	PaginationOpt *model.PaginationOptionInput
}
