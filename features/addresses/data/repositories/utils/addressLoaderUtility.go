package addressdomainrepositoryutilities

import (
	"context"

	googlemapcoreoperationinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/operations"
	addressdomainrepositorytypes "github.com/horeekaa/backend/features/addresses/domain/repositories/types"
	addressdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
	"googlemaps.github.io/maps"
)

type addressLoader struct {
	gMapOperation googlemapcoreoperationinterfaces.GoogleMapBasicOperation
}

func NewAddressLoader(
	gMapOperation googlemapcoreoperationinterfaces.GoogleMapBasicOperation,
) (addressdomainrepositoryutilityinterfaces.AddressLoader, error) {
	return &addressLoader{
		gMapOperation,
	}, nil
}

func (addrLoader *addressLoader) Execute(
	input *addressdomainrepositorytypes.LatLngGeocode,
	resolvedGeocoding *model.ResolvedGeocodingInput,
) (bool, error) {
	geocodeResults, err := addrLoader.gMapOperation.ReverseGeocode(
		context.Background(),
		&maps.GeocodingRequest{
			LatLng: &maps.LatLng{
				Lat: input.Latitude,
				Lng: input.Longitude,
			},
		},
	)
	if err != nil {
		return false, err
	}

	for _, addrComponent := range geocodeResults[0].AddressComponents {
		if funk.Contains(
			addrComponent.Types,
			func(s string) bool {
				return s == "administrative_area_level_3"
			},
		) {
			*resolvedGeocoding.District = addrComponent.LongName
		}

		if funk.Contains(
			addrComponent.Types,
			func(s string) bool {
				return s == "administrative_area_level_2"
			},
		) {
			*resolvedGeocoding.Municipality = addrComponent.LongName
		}

		if funk.Contains(
			addrComponent.Types,
			func(s string) bool {
				return s == "administrative_area_level_1"
			},
		) {
			*resolvedGeocoding.Province = addrComponent.LongName
		}

		if funk.Contains(
			addrComponent.Types,
			func(s string) bool {
				return s == "postal_code"
			},
		) {
			*resolvedGeocoding.ZipCode = addrComponent.LongName
		}

		if funk.Contains(
			addrComponent.Types,
			func(s string) bool {
				return s == "street_address"
			},
		) {
			*resolvedGeocoding.StreetName = addrComponent.LongName
		}
	}

	return true, nil
}
