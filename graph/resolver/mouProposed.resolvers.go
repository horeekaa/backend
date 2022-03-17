package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	mouitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/mouItems/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mouProposedResolver) Items(ctx context.Context, obj *model.MouProposed) ([]*model.MouItem, error) {
	var getMouItemUsecase mouitempresentationusecaseinterfaces.GetMouItemUsecase
	container.Make(&getMouItemUsecase)

	mouItems := []*model.MouItem{}
	if obj.Items != nil {
		for _, item := range obj.Items {
			mouItem, err := getMouItemUsecase.Execute(
				&model.MouItemFilterFields{
					ID: &item.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			mouItems = append(mouItems, mouItem)
		}
	}
	return mouItems, nil
}

func (r *mouProposedResolver) SubmittingAccount(ctx context.Context, obj *model.MouProposed) (*model.Account, error) {
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

func (r *mouProposedResolver) RecentApprovingAccount(ctx context.Context, obj *model.MouProposed) (*model.Account, error) {
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

func (r *mouProposedResolver) RecentLog(ctx context.Context, obj *model.MouProposed) (*model.Logging, error) {
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

// MouProposed returns generated.MouProposedResolver implementation.
func (r *Resolver) MouProposed() generated.MouProposedResolver { return &mouProposedResolver{r} }

type mouProposedResolver struct{ *Resolver }
