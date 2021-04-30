package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *loggingResolver) CreatedByAccount(ctx context.Context, obj *model.Logging) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

// Logging returns generated.LoggingResolver implementation.
func (r *Resolver) Logging() generated.LoggingResolver { return &loggingResolver{r} }

type loggingResolver struct{ *Resolver }
