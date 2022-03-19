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
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *supplyOrderProposedResolver) Items(ctx context.Context, obj *model.SupplyOrderProposed) ([]*model.SupplyOrderItem, error) {
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

func (r *supplyOrderProposedResolver) SubmittingAccount(ctx context.Context, obj *model.SupplyOrderProposed) (*model.Account, error) {
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

func (r *supplyOrderProposedResolver) RecentApprovingAccount(ctx context.Context, obj *model.SupplyOrderProposed) (*model.Account, error) {
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

func (r *supplyOrderProposedResolver) RecentLog(ctx context.Context, obj *model.SupplyOrderProposed) (*model.Logging, error) {
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

// SupplyOrderProposed returns generated.SupplyOrderProposedResolver implementation.
func (r *Resolver) SupplyOrderProposed() generated.SupplyOrderProposedResolver {
	return &supplyOrderProposedResolver{r}
}

type supplyOrderProposedResolver struct{ *Resolver }
