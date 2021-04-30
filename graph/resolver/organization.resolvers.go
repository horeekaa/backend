package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *organizationResolver) SubmittingAccount(ctx context.Context, obj *model.Organization) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *organizationResolver) ApprovingAccount(ctx context.Context, obj *model.Organization) (*model.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *organizationResolver) PreviousEntity(ctx context.Context, obj *model.Organization) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *organizationResolver) CorrespondingLog(ctx context.Context, obj *model.Organization) (*model.Logging, error) {
	panic(fmt.Errorf("not implemented"))
}

// Organization returns generated.OrganizationResolver implementation.
func (r *Resolver) Organization() generated.OrganizationResolver { return &organizationResolver{r} }

type organizationResolver struct{ *Resolver }
