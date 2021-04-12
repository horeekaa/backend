package graphresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/horeekaa/backend/graph/generated"
	"github.com/horeekaa/backend/model"
)

func (r *memberAccessRefResolver) SubmittingPerson(ctx context.Context, obj *model.MemberAccessRef) (*model.Person, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *memberAccessRefResolver) ApprovingPerson(ctx context.Context, obj *model.MemberAccessRef) (*model.Person, error) {
	panic(fmt.Errorf("not implemented"))
}

// MemberAccessRef returns generated.MemberAccessRefResolver implementation.
func (r *Resolver) MemberAccessRef() generated.MemberAccessRefResolver {
	return &memberAccessRefResolver{r}
}

type memberAccessRefResolver struct{ *Resolver }
