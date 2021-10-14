package addressdomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	addressdomainrepositorytypes "github.com/horeekaa/backend/features/addresses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type AddressLoader interface {
	Execute(
		operationOptions *mongodbcoretypes.OperationOptions,
		input *addressdomainrepositorytypes.LatLngGeocode,
		resolvedGeocoding *model.ResolvedGeocodingInput,
		addressRegion *model.AddressRegionGroupForAddressInput,
	) (bool, error)
}
