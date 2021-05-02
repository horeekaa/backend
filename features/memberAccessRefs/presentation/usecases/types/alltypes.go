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

type UpdateMemberAccessRefUsecaseInput struct {
	AuthHeader            string
	Context               context.Context
	UpdateMemberAccessRef *model.UpdateMemberAccessRef
}

type GetAllMemberAccessRefUsecaseInput struct {
	AuthHeader    string
	Context       context.Context
	FilterFields  *model.MemberAccessRefFilterFields
	PaginationOps *model.PaginationOptionInput
}
