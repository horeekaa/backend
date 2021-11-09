package purchaseordertosupplydomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllPurchaseOrderToSupplyInput struct {
	FilterFields  *model.PurchaseOrderToSupplyFilterFields
	PaginationOpt *model.PaginationOptionInput
}
