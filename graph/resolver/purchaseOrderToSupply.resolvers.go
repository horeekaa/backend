package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	purchaseorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases"
	purchaseorderitempresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases/types"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
	purchaseordertosupplypresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases/types"
	supplyorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases"
	supplyorderitempresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *purchaseOrderToSupplyResolver) PurchaseOrderItems(ctx context.Context, obj *model.PurchaseOrderToSupply) ([]*model.PurchaseOrderItem, error) {
	var getAllPurchaseOrderItemUsecase purchaseorderitempresentationusecaseinterfaces.GetAllPurchaseOrderItemUsecase
	container.Make(&getAllPurchaseOrderItemUsecase)

	if obj.PurchaseOrderItems != nil {
		purchaseOrderItems, err := getAllPurchaseOrderItemUsecase.Execute(
			purchaseorderitempresentationusecasetypes.GetAllPurchaseOrderItemUsecaseInput{
				Context: ctx,
				FilterFields: &model.PurchaseOrderItemFilterFields{
					ID: &model.ObjectIDFilterField{
						Operation: model.ObjectIDOperationIn,
						Values: funk.Map(
							obj.PurchaseOrderItems,
							func(poItem *model.PurchaseOrderItem) *primitive.ObjectID {
								return &poItem.ID
							},
						).([]*primitive.ObjectID),
					},
				},
				PaginationOps: &model.PaginationOptionInput{
					QueryLimit: func(i int) *int { return &i }(999),
				},
			},
		)
		if err != nil {
			return []*model.PurchaseOrderItem{}, nil
		}
		return purchaseOrderItems, nil
	}
	return []*model.PurchaseOrderItem{}, nil
}

func (r *purchaseOrderToSupplyResolver) SupplyOrderItems(ctx context.Context, obj *model.PurchaseOrderToSupply) ([]*model.SupplyOrderItem, error) {
	var getAllSupplyOrderItemUsecase supplyorderitempresentationusecaseinterfaces.GetAllSupplyOrderItemUsecase
	container.Make(&getAllSupplyOrderItemUsecase)

	if obj.SupplyOrderItems != nil {
		supplyOrderItems, err := getAllSupplyOrderItemUsecase.Execute(
			supplyorderitempresentationusecasetypes.GetAllSupplyOrderItemUsecaseInput{
				Context: ctx,
				FilterFields: &model.SupplyOrderItemFilterFields{
					ID: &model.ObjectIDFilterField{
						Operation: model.ObjectIDOperationIn,
						Values: funk.Map(
							obj.SupplyOrderItems,
							func(soItem *model.SupplyOrderItem) *primitive.ObjectID {
								return &soItem.ID
							},
						).([]*primitive.ObjectID),
					},
				},
				PaginationOps: &model.PaginationOptionInput{
					QueryLimit: func(i int) *int { return &i }(999),
				},
			},
		)
		if err != nil {
			return []*model.SupplyOrderItem{}, nil
		}
		return supplyOrderItems, nil
	}
	return []*model.SupplyOrderItem{}, nil
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
