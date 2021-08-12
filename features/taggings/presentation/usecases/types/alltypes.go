package taggingpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type BulkCreateTaggingUsecaseInput struct {
	Context           context.Context
	BulkCreateTagging *model.BulkCreateTagging
}

type BulkUpdateTaggingUsecaseInput struct {
	Context           context.Context
	BulkUpdateTagging *model.BulkUpdateTagging
}

type GetAllTaggingUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.TaggingFilterFields
	PaginationOps *model.PaginationOptionInput
}
