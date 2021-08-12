package tagdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllTagInput struct {
	FilterFields  *model.TagFilterFields
	PaginationOpt *model.PaginationOptionInput
}
