package memberaccesspresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessUsecaseInput struct {
	Context            context.Context
	CreateMemberAccess *model.CreateMemberAccess
}

type UpdateMemberAccessUsecaseInput struct {
	Context            context.Context
	UpdateMemberAccess *model.UpdateMemberAccess
}

type GetAllMemberAccessUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.MemberAccessFilterFields
	PaginationOps *model.PaginationOptionInput
}
