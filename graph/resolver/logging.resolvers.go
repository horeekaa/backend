package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/golobby/container/v2"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *loggingResolver) CreatedByAccount(ctx context.Context, obj *model.Logging) (*model.Account, error) {
	var getAccountUsecase accountpresentationusecaseinterfaces.GetAccountUsecase
	container.Make(&getAccountUsecase)
	return getAccountUsecase.Execute(
		&model.AccountFilterFields{
			ID: &obj.CreatedByAccount.ID,
		},
	)
}

// Logging returns generated.LoggingResolver implementation.
func (r *Resolver) Logging() generated.LoggingResolver { return &loggingResolver{r} }

type loggingResolver struct{ *Resolver }
