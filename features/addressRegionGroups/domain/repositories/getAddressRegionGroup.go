package addressregiongroupdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetAddressRegionGroupRepository interface {
	Execute(filterFields *model.AddressRegionGroupFilterFields) (*model.AddressRegionGroup, error)
}
