package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *accountResolver) Person(ctx context.Context, obj *model.Account) (*model.Person, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *memberAccessResolver) Account(ctx context.Context, obj *model.MemberAccess) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *memberAccessResolver) Organization(ctx context.Context, obj *model.MemberAccess) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *memberAccessResolver) DefaultAccess(ctx context.Context, obj *model.MemberAccess) (*model.MemberAccessRef, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Login(ctx context.Context, deviceToken *string) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Logout(ctx context.Context, deviceToken *string) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateOrganization(ctx context.Context, createOrganization model.CreateOrganization) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context) (*model.PersonDetailed, error) {
	panic(fmt.Errorf("not implemented"))
}

// Account returns generated.AccountResolver implementation.
func (r *Resolver) Account() generated.AccountResolver { return &accountResolver{r} }

// MemberAccess returns generated.MemberAccessResolver implementation.
func (r *Resolver) MemberAccess() generated.MemberAccessResolver { return &memberAccessResolver{r} }

type accountResolver struct{ *Resolver }
type memberAccessResolver struct{ *Resolver }
