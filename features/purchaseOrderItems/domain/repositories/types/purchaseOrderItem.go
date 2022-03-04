package purchaseorderitemdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllPurchaseOrderItemInput struct {
	FilterFields  *model.PurchaseOrderItemFilterFields
	PaginationOpt *model.PaginationOptionInput
}
