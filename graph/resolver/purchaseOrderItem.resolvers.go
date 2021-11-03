package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *purchaseOrderItemResolver) PurchaseOrder(ctx context.Context, obj *model.PurchaseOrderItem) (*model.PurchaseOrder, error) {
	var getPurchaseOrderUsecase purchaseorderpresentationusecaseinterfaces.GetPurchaseOrderUsecase
	container.Make(&getPurchaseOrderUsecase)
	return getPurchaseOrderUsecase.Execute(
		&model.PurchaseOrderFilterFields{
			ID: &obj.PurchaseOrder.ID,
		},
	)
}

func (r *purchaseOrderItemResolver) PurchaseOrderToSupply(ctx context.Context, obj *model.PurchaseOrderItem) (*model.PurchaseOrderToSupply, error) {
	var getPurchaseOrderToSupplyUsecase purchaseordertosupplypresentationusecaseinterfaces.GetPurchaseOrderToSupplyUsecase
	container.Make(&getPurchaseOrderToSupplyUsecase)
	return getPurchaseOrderToSupplyUsecase.Execute(
		&model.PurchaseOrderToSupplyFilterFields{
			ID: &obj.PurchaseOrderToSupply.ID,
		},
	)
}

// PurchaseOrderItem returns generated.PurchaseOrderItemResolver implementation.
func (r *Resolver) PurchaseOrderItem() generated.PurchaseOrderItemResolver {
	return &purchaseOrderItemResolver{r}
}

type purchaseOrderItemResolver struct{ *Resolver }
