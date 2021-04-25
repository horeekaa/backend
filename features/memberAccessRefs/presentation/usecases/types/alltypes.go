package memberaccessrefpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessRefUsecaseInput struct {
	AuthHeader            string
	Context               context.Context
	CreateMemberAccessRef *model.CreateMemberAccessRef
}
