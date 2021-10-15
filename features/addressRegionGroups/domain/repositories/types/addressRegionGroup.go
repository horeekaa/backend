package addressregiongroupdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllAddressRegionGroupInput struct {
	FilterFields  *model.AddressRegionGroupFilterFields
	PaginationOpt *model.PaginationOptionInput
}
