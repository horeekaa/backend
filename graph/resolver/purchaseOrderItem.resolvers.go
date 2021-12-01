package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
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

	var filterFields *model.PurchaseOrderToSupplyFilterFields
	if obj.PurchaseOrderToSupply != nil {
		filterFields = &model.PurchaseOrderToSupplyFilterFields{}
		filterFields.ID = &obj.PurchaseOrderToSupply.ID
	}
	return getPurchaseOrderToSupplyUsecase.Execute(
		filterFields,
	)
}

func (r *purchaseOrderItemResolver) SubmittingAccount(ctx context.Context, obj *model.PurchaseOrderItem) (*model.Account, error) {
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

func (r *purchaseOrderItemResolver) RecentApprovingAccount(ctx context.Context, obj *model.PurchaseOrderItem) (*model.Account, error) {
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

func (r *purchaseOrderItemResolver) RecentLog(ctx context.Context, obj *model.PurchaseOrderItem) (*model.Logging, error) {
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

// PurchaseOrderItem returns generated.PurchaseOrderItemResolver implementation.
func (r *Resolver) PurchaseOrderItem() generated.PurchaseOrderItemResolver {
	return &purchaseOrderItemResolver{r}
}

type purchaseOrderItemResolver struct{ *Resolver }
