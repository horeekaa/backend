package moupresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreateMouUsecaseInput struct {
	Context   context.Context
	CreateMou *model.CreateMou
}

type UpdateMouUsecaseInput struct {
	Context   context.Context
	UpdateMou *model.UpdateMou
}

type GetAllMouUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.MouFilterFields
	PaginationOps *model.PaginationOptionInput
}
