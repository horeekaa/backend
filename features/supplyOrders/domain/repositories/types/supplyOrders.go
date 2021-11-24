package supplyorderdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllSupplyOrderInput struct {
	FilterFields  *model.SupplyOrderFilterFields
	PaginationOpt *model.PaginationOptionInput
}
