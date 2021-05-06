package memberaccessrefpresentationusecasetypes

import (
	"context"

	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessRefUsecaseInput struct {
	User                  *authenticationcoremodels.AuthUserWrap
	Context               context.Context
	CreateMemberAccessRef *model.CreateMemberAccessRef
}

type UpdateMemberAccessRefUsecaseInput struct {
	User                  *authenticationcoremodels.AuthUserWrap
	Context               context.Context
	UpdateMemberAccessRef *model.UpdateMemberAccessRef
}

type GetAllMemberAccessRefUsecaseInput struct {
	User          *authenticationcoremodels.AuthUserWrap
	Context       context.Context
	FilterFields  *model.MemberAccessRefFilterFields
	PaginationOps *model.PaginationOptionInput
}
