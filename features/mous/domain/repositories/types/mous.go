package moudomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllMouInput struct {
	FilterFields  *model.MouFilterFields
	PaginationOpt *model.PaginationOptionInput
}
