package taggingdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllTaggingInput struct {
	FilterFields  *model.TaggingFilterFields
	PaginationOpt *model.PaginationOptionInput
}
