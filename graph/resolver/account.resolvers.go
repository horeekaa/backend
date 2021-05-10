package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	container "github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *accountResolver) Person(ctx context.Context, obj *model.Account) (*model.Person, error) {
	var getPersonDataFromAccountUsecase accountpresentationusecaseinterfaces.GetPersonDataFromAccountUsecase
	container.Make(&getPersonDataFromAccountUsecase)
	return getPersonDataFromAccountUsecase.Execute(
		accountpresentationusecasetypes.GetPersonDataFromAccountInput{
			Account:         obj,
			Context:         ctx,
			ViewProfileMode: false,
		},
	)
}

func (r *mutationResolver) Login(ctx context.Context, deviceToken *string) (*model.Account, error) {
	var loginUsecase accountpresentationusecaseinterfaces.LoginUsecase
	container.Make(&loginUsecase)
	return loginUsecase.Execute(
		accountpresentationusecasetypes.LoginUsecaseInput{
			DeviceToken: *deviceToken,
			Context:     ctx,
		},
	)
}

func (r *mutationResolver) Logout(ctx context.Context, deviceToken *string) (*model.Account, error) {
	var logoutUsecase accountpresentationusecaseinterfaces.LogoutUsecase
	container.Make(&logoutUsecase)
	return logoutUsecase.Execute(
		accountpresentationusecasetypes.LogoutUsecaseInput{
			DeviceToken: *deviceToken,
			Context:     ctx,
		},
	)
}

func (r *queryResolver) Person(ctx context.Context, account *model.ObjectIDOnly) (*model.Person, error) {
	var getPersonDataFromAccountUsecase accountpresentationusecaseinterfaces.GetPersonDataFromAccountUsecase
	container.Make(&getPersonDataFromAccountUsecase)
	return getPersonDataFromAccountUsecase.Execute(
		accountpresentationusecasetypes.GetPersonDataFromAccountInput{
			Account:         &model.Account{ID: *account.ID},
			Context:         ctx,
			ViewProfileMode: true,
		},
	)
}

// Account returns generated.AccountResolver implementation.
func (r *Resolver) Account() generated.AccountResolver { return &accountResolver{r} }

type accountResolver struct{ *Resolver }
