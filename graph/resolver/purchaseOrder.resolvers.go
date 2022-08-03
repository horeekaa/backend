package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	purchaseorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
	purchaseorderpresentationusecasetypes "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mutationResolver) CreatePurchaseOrder(ctx context.Context, createPurchaseOrder model.CreatePurchaseOrder) ([]*model.PurchaseOrder, error) {
	var createPurchaseOrderUsecase purchaseorderpresentationusecaseinterfaces.CreatePurchaseOrderUsecase
	container.Make(&createPurchaseOrderUsecase)
	return createPurchaseOrderUsecase.Execute(
		purchaseorderpresentationusecasetypes.CreatePurchaseOrderUsecaseInput{
			Context:             ctx,
			CreatePurchaseOrder: &createPurchaseOrder,
		},
	)
}

func (r *mutationResolver) UpdatePurchaseOrder(ctx context.Context, updatePurchaseOrder model.UpdatePurchaseOrder) (*model.PurchaseOrder, error) {
	var updatePurchaseOrderUsecase purchaseorderpresentationusecaseinterfaces.UpdatePurchaseOrderUsecase
	container.Make(&updatePurchaseOrderUsecase)
	return updatePurchaseOrderUsecase.Execute(
		purchaseorderpresentationusecasetypes.UpdatePurchaseOrderUsecaseInput{
			Context:             ctx,
			UpdatePurchaseOrder: &updatePurchaseOrder,
		},
	)
}

func (r *purchaseOrderResolver) Items(ctx context.Context, obj *model.PurchaseOrder) ([]*model.PurchaseOrderItem, error) {
	var getPurchaseOrderItemUsecase purchaseorderitempresentationusecaseinterfaces.GetPurchaseOrderItemUsecase
	container.Make(&getPurchaseOrderItemUsecase)

	purchaseOrderItems := []*model.PurchaseOrderItem{}
	if obj.Items != nil {
		for _, item := range obj.Items {
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

func (r *purchaseOrderResolver) SubmittingAccount(ctx context.Context, obj *model.PurchaseOrder) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)
	return getAccountUsecase.Execute(
		accountpresentationusecasetypes.GetAccountInput{
			FilterFields: &model.AccountFilterFields{
				ID: &obj.SubmittingAccount.ID,
			},
		},
	)
}

func (r *purchaseOrderResolver) RecentApprovingAccount(ctx context.Context, obj *model.PurchaseOrder) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)

	var filterFields *model.AccountFilterFields
	if obj.RecentApprovingAccount != nil {
		filterFields = &model.AccountFilterFields{}
		filterFields.ID = &obj.RecentApprovingAccount.ID
	}
	return getAccountUsecase.Execute(
		accountpresentationusecasetypes.GetAccountInput{
			FilterFields: filterFields,
		},
	)
}

func (r *purchaseOrderResolver) RecentLog(ctx context.Context, obj *model.PurchaseOrder) (*model.Logging, error) {
	var getLoggingUsecase loggingpresentationusecaseinterfaces.GetLoggingUsecase
	container.Make(&getLoggingUsecase)

	var filterFields *model.LoggingFilterFields
	if obj.RecentLog != nil {
		filterFields = &model.LoggingFilterFields{}
		filterFields.ID = &obj.RecentLog.ID
	}
	return getLoggingUsecase.Execute(
		filterFields,
	)
}

func (r *queryResolver) PurchaseOrders(ctx context.Context, filterFields model.PurchaseOrderFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.PurchaseOrder, error) {
	var getPurchaseOrdersUsecase purchaseorderpresentationusecaseinterfaces.GetAllPurchaseOrderUsecase
	container.Make(&getPurchaseOrdersUsecase)
	return getPurchaseOrdersUsecase.Execute(
		purchaseorderpresentationusecasetypes.GetAllPurchaseOrderUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

// PurchaseOrder returns generated.PurchaseOrderResolver implementation.
func (r *Resolver) PurchaseOrder() generated.PurchaseOrderResolver { return &purchaseOrderResolver{r} }

type purchaseOrderResolver struct{ *Resolver }
