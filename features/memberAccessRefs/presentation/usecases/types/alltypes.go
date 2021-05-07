package memberaccessrefpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessRefUsecaseInput struct {
	Context               context.Context
	CreateMemberAccessRef *model.CreateMemberAccessRef
}

type UpdateMemberAccessRefUsecaseInput struct {
	Context               context.Context
	UpdateMemberAccessRef *model.UpdateMemberAccessRef
}

type GetAllMemberAccessRefUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.MemberAccessRefFilterFields
	PaginationOps *model.PaginationOptionInput
}
