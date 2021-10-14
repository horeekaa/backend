package addressdomainrepositoryutilities

import (
	"context"
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	googlemapcoreoperationinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/operations"
	databaseaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	addressdomainrepositorytypes "github.com/horeekaa/backend/features/addresses/domain/repositories/types"
	addressdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
	"googlemaps.github.io/maps"
)

type addressLoader struct {
	gMapOperation           googlemapcoreoperationinterfaces.GoogleMapBasicOperation
	addressRegionDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource
}

func NewAddressLoader(
	gMapOperation googlemapcoreoperationinterfaces.GoogleMapBasicOperation,
	addressRegionDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource,
) (addressdomainrepositoryutilityinterfaces.AddressLoader, error) {
	return &addressLoader{
		gMapOperation,
		addressRegionDataSource,
	}, nil
}

func (addrLoader *addressLoader) Execute(
	operationOptions *mongodbcoretypes.OperationOptions,
	input *addressdomainrepositorytypes.LatLngGeocode,
	resolvedGeocoding *model.ResolvedGeocodingInput,
	addressRegionToLoad *model.AddressRegionGroupForAddressInput,
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

	addressRegion, err := addrLoader.addressRegionDataSource.GetMongoDataSource().FindOne(
		map[string]interface{}{
			"cities": map[string]interface{}{
				"$elemMatch": map[string]interface{}{
					"$regex":   *resolvedGeocoding.Municipality,
					"$options": "i",
				},
			},
		},
		operationOptions,
	)
	if err != nil {
		return false, err
	}

	if addressRegion != nil {
		jsonTemp, _ := json.Marshal(addressRegion)
		json.Unmarshal(jsonTemp, addressRegionToLoad)
	}

	return true, nil
}
