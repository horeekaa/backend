package addressregiongroupdomainrepositoryinterfaces

import (
	addressregiongroupdomainrepositorytypes "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllAddressRegionGroupRepository interface {
	Execute(filterFields addressregiongroupdomainrepositorytypes.GetAllAddressRegionGroupInput) ([]*model.AddressRegionGroup, error)
}
