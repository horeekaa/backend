package purchaseorderpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreatePurchaseOrderUsecaseInput struct {
	Context             context.Context
	CreatePurchaseOrder *model.CreatePurchaseOrder
}

type UpdatePurchaseOrderUsecaseInput struct {
	Context             context.Context
	UpdatePurchaseOrder *model.UpdatePurchaseOrder
}

type GetAllPurchaseOrderUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.PurchaseOrderFilterFields
	PaginationOps *model.PaginationOptionInput
}
