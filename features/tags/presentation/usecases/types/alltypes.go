package tagpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreateTagUsecaseInput struct {
	Context   context.Context
	CreateTag *model.CreateTag
}

type UpdateTagUsecaseInput struct {
	Context   context.Context
	UpdateTag *model.UpdateTag
}

type GetAllTagUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.TagFilterFields
	PaginationOps *model.PaginationOptionInput
}
