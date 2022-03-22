package mouitemdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllMouItemInput struct {
	FilterFields  *model.MouItemFilterFields
	PaginationOpt *model.PaginationOptionInput
}
