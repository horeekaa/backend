package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	moupresentationusecaseinterfaces "github.com/horeekaa/backend/features/mous/presentation/usecases"
	moupresentationusecasetypes "github.com/horeekaa/backend/features/mous/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *mouResolver) SubmittingAccount(ctx context.Context, obj *model.Mou) (*model.Account, error) {
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

func (r *mouResolver) RecentApprovingAccount(ctx context.Context, obj *model.Mou) (*model.Account, error) {
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

func (r *mouResolver) RecentLog(ctx context.Context, obj *model.Mou) (*model.Logging, error) {
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

func (r *mutationResolver) CreateMou(ctx context.Context, createMou model.CreateMou) (*model.Mou, error) {
	var createMouUsecase moupresentationusecaseinterfaces.CreateMouUsecase
	container.Make(&createMouUsecase)
	return createMouUsecase.Execute(
		moupresentationusecasetypes.CreateMouUsecaseInput{
			Context:   ctx,
			CreateMou: &createMou,
		},
	)
}

func (r *mutationResolver) UpdateMou(ctx context.Context, updateMou model.UpdateMou) (*model.Mou, error) {
	var updateMouUsecase moupresentationusecaseinterfaces.UpdateMouUsecase
	container.Make(&updateMouUsecase)
	return updateMouUsecase.Execute(
		moupresentationusecasetypes.UpdateMouUsecaseInput{
			Context:   ctx,
			UpdateMou: &updateMou,
		},
	)
}

func (r *queryResolver) Mous(ctx context.Context, filterFields model.MouFilterFields, paginationOpt *model.PaginationOptionInput) ([]*model.Mou, error) {
	var getMousUsecass moupresentationusecaseinterfaces.GetAllMouUsecase
	container.Make(&getMousUsecass)
	return getMousUsecass.Execute(
		moupresentationusecasetypes.GetAllMouUsecaseInput{
			Context:       ctx,
			FilterFields:  &filterFields,
			PaginationOps: paginationOpt,
		},
	)
}

// Mou returns generated.MouResolver implementation.
func (r *Resolver) Mou() generated.MouResolver { return &mouResolver{r} }

type mouResolver struct{ *Resolver }
