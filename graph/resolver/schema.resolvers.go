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

func (r *mutationResolver) Login(ctx context.Context, deviceToken *string) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Logout(ctx context.Context, deviceToken *string) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateOrganization(ctx context.Context, newOrganization model.NewOrganization) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *organizationMembershipResolver) DefaultAccess(ctx context.Context, obj *model.OrganizationMembership) (*model.MemberAccess, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *personResolver) OrganizationMembership(ctx context.Context, obj *model.Person) (*model.OrganizationMembership, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context) (*model.PersonDetailed, error) {
	panic(fmt.Errorf("not implemented"))
}

// Account returns generated.AccountResolver implementation.
func (r *Resolver) Account() generated.AccountResolver { return &accountResolver{r} }

// OrganizationMembership returns generated.OrganizationMembershipResolver implementation.
func (r *Resolver) OrganizationMembership() generated.OrganizationMembershipResolver {
	return &organizationMembershipResolver{r}
}

// Person returns generated.PersonResolver implementation.
func (r *Resolver) Person() generated.PersonResolver { return &personResolver{r} }

type accountResolver struct{ *Resolver }
type organizationMembershipResolver struct{ *Resolver }
type personResolver struct{ *Resolver }
