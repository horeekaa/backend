package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	purchaseorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
	purchaseordertosupplypresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *purchaseOrderToSupplyResolver) PurchaseOrderItems(ctx context.Context, obj *model.PurchaseOrderToSupply) ([]*model.PurchaseOrderItem, error) {
	var getPurchaseOrderItemUsecase purchaseorderitempresentationusecaseinterfaces.GetPurchaseOrderItemUsecase
	container.Make(&getPurchaseOrderItemUsecase)

	purchaseOrderItems := []*model.PurchaseOrderItem{}
	if obj.PurchaseOrderItems != nil {
		for _, item := range obj.PurchaseOrderItems {
			purchaseOrderItem, err := getPurchaseOrderItemUsecase.Execute(
				&model.PurchaseOrderItemFilterFields{
					ID: &item.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			purchaseOrderItems = append(purchaseOrderItems, purchaseOrderItem)
		}
	}
	return purchaseOrderItems, nil
}

func (r *queryResolver) PurchaseOrdersToSupply(ctx context.Context, filterFields model.PurchaseOrderToSupplyFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.PurchaseOrderToSupply, error) {
	var getPurchaseOrdersToSupplyUsecase purchaseordertosupplypresentationusecaseinterfaces.GetAllPurchaseOrderToSupplyUsecase
	container.Make(&getPurchaseOrdersToSupplyUsecase)
	return getPurchaseOrdersToSupplyUsecase.Execute(
		purchaseordertosupplypresentationusecasetypes.GetAllPurchaseOrderToSupplyUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

// PurchaseOrderToSupply returns generated.PurchaseOrderToSupplyResolver implementation.
func (r *Resolver) PurchaseOrderToSupply() generated.PurchaseOrderToSupplyResolver {
	return &purchaseOrderToSupplyResolver{r}
}

type purchaseOrderToSupplyResolver struct{ *Resolver }
