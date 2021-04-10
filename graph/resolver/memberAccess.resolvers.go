package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *memberAccessResolver) Account(ctx context.Context, obj *model.MemberAccess) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *memberAccessResolver) Organization(ctx context.Context, obj *model.MemberAccess) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *memberAccessResolver) DefaultAccess(ctx context.Context, obj *model.MemberAccess) (*model.MemberAccessRef, error) {
	panic(fmt.Errorf("not implemented"))
}

// MemberAccess returns generated.MemberAccessResolver implementation.
func (r *Resolver) MemberAccess() generated.MemberAccessResolver { return &memberAccessResolver{r} }

type memberAccessResolver struct{ *Resolver }
