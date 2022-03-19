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
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *purchaseOrderProposedResolver) Items(ctx context.Context, obj *model.PurchaseOrderProposed) ([]*model.PurchaseOrderItem, error) {
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

func (r *purchaseOrderProposedResolver) SubmittingAccount(ctx context.Context, obj *model.PurchaseOrderProposed) (*model.Account, error) {
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

func (r *purchaseOrderProposedResolver) RecentApprovingAccount(ctx context.Context, obj *model.PurchaseOrderProposed) (*model.Account, error) {
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

func (r *purchaseOrderProposedResolver) RecentLog(ctx context.Context, obj *model.PurchaseOrderProposed) (*model.Logging, error) {
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

// PurchaseOrderProposed returns generated.PurchaseOrderProposedResolver implementation.
func (r *Resolver) PurchaseOrderProposed() generated.PurchaseOrderProposedResolver {
	return &purchaseOrderProposedResolver{r}
}

type purchaseOrderProposedResolver struct{ *Resolver }
