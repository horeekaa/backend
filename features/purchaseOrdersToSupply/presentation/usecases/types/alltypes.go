package purchaseordertosupplypresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type GetAllPurchaseOrderToSupplyUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.PurchaseOrderToSupplyFilterFields
	PaginationOps *model.PaginationOptionInput
}
