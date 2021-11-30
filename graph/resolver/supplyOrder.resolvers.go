package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	supplyorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases"
	supplyorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases"
	supplyorderpresentationusecasetypes "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mutationResolver) CreateSupplyOrder(ctx context.Context, createSupplyOrder model.CreateSupplyOrder) ([]*model.SupplyOrder, error) {
	var createSupplyOrderUsecase supplyorderpresentationusecaseinterfaces.CreateSupplyOrderUsecase
	container.Make(&createSupplyOrderUsecase)
	return createSupplyOrderUsecase.Execute(
		supplyorderpresentationusecasetypes.CreateSupplyOrderUsecaseInput{
			Context:           ctx,
			CreateSupplyOrder: &createSupplyOrder,
		},
	)
}

func (r *mutationResolver) UpdateSupplyOrder(ctx context.Context, updateSupplyOrder model.UpdateSupplyOrder) (*model.SupplyOrder, error) {
	var updateSupplyOrderUsecase supplyorderpresentationusecaseinterfaces.UpdateSupplyOrderUsecase
	container.Make(&updateSupplyOrderUsecase)
	return updateSupplyOrderUsecase.Execute(
		supplyorderpresentationusecasetypes.UpdateSupplyOrderUsecaseInput{
			Context:           ctx,
			UpdateSupplyOrder: &updateSupplyOrder,
		},
	)
}

func (r *queryResolver) SupplyOrders(ctx context.Context, filterFields model.SupplyOrderFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.SupplyOrder, error) {
	var getSupplyOrdersUsecase supplyorderpresentationusecaseinterfaces.GetAllSupplyOrderUsecase
	container.Make(&getSupplyOrdersUsecase)
	return getSupplyOrdersUsecase.Execute(
		supplyorderpresentationusecasetypes.GetAllSupplyOrderUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

func (r *supplyOrderResolver) Items(ctx context.Context, obj *model.SupplyOrder) ([]*model.SupplyOrderItem, error) {
	var getSupplyOrderItemUsecase supplyorderitempresentationusecaseinterfaces.GetSupplyOrderItemUsecase
	container.Make(&getSupplyOrderItemUsecase)

	supplyOrderItems := []*model.SupplyOrderItem{}
	if obj.Items != nil {
		for _, item := range obj.Items {
			supplyOrderItem, err := getSupplyOrderItemUsecase.Execute(
				&model.SupplyOrderItemFilterFields{
					ID: &item.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			supplyOrderItems = append(supplyOrderItems, supplyOrderItem)
		}
	}
	return supplyOrderItems, nil
}

func (r *supplyOrderResolver) SubmittingAccount(ctx context.Context, obj *model.SupplyOrder) (*model.Account, error) {
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

func (r *supplyOrderResolver) RecentApprovingAccount(ctx context.Context, obj *model.SupplyOrder) (*model.Account, error) {
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

func (r *supplyOrderResolver) RecentLog(ctx context.Context, obj *model.SupplyOrder) (*model.Logging, error) {
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

// SupplyOrder returns generated.SupplyOrderResolver implementation.
func (r *Resolver) SupplyOrder() generated.SupplyOrderResolver { return &supplyOrderResolver{r} }

type supplyOrderResolver struct{ *Resolver }
