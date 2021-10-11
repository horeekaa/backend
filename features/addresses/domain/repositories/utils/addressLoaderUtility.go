package addressdomainrepositoryutilityinterfaces

import (
	addressdomainrepositorytypes "github.com/horeekaa/backend/features/addresses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type AddressLoader interface {
	Execute(
		input *addressdomainrepositorytypes.LatLngGeocode,
		resolvedGeocoding *model.ResolvedGeocodingInput,
	) (bool, error)
}
