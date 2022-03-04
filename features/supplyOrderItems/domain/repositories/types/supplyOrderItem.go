package supplyorderitemdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllSupplyOrderItemInput struct {
	FilterFields  *model.SupplyOrderItemFilterFields
	PaginationOpt *model.PaginationOptionInput
}
