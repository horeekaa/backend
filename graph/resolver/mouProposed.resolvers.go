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
	mouitempresentationusecasetypes "github.com/horeekaa/backend/features/mouItems/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mouProposedResolver) Items(ctx context.Context, obj *model.MouProposed) ([]*model.MouItem, error) {
	var getAllMouItemUsecase mouitempresentationusecaseinterfaces.GetAllMouItemUsecase
	container.Make(&getAllMouItemUsecase)

	if obj.Items != nil {
		mouItems, err := getAllMouItemUsecase.Execute(
			mouitempresentationusecasetypes.GetAllMouItemUsecaseInput{
				Context: ctx,
				FilterFields: &model.MouItemFilterFields{
					ID: &model.ObjectIDFilterField{
						Operation: model.ObjectIDOperationIn,
						Values: funk.Map(
							obj.Items,
							func(item *model.MouItem) *primitive.ObjectID {
								return &item.ID
							},
						).([]*primitive.ObjectID),
					},
				},
				PaginationOps: &model.PaginationOptionInput{},
			},
		)
		if err != nil {
			return nil, err
		}
		return mouItems, nil
	}
	return []*model.MouItem{}, nil
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
