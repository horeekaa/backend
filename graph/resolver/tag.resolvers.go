package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *tagResolver) Photos(ctx context.Context, obj *model.Tag) ([]*model.DescriptivePhoto, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *tagResolver) SubmittingAccount(ctx context.Context, obj *model.Tag) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *tagResolver) RecentApprovingAccount(ctx context.Context, obj *model.Tag) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *tagResolver) RecentLog(ctx context.Context, obj *model.Tag) (*model.Logging, error) {
	panic(fmt.Errorf("not implemented"))
}

// Tag returns generated.TagResolver implementation.
func (r *Resolver) Tag() generated.TagResolver { return &tagResolver{r} }

type tagResolver struct{ *Resolver }
