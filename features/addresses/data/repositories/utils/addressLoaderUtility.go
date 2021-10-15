package addressdomainrepositoryutilities

import (
	"context"
	"encoding/json"
	"strings"

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
	resolvedGeocodingToLoad *model.ResolvedGeocodingInput,
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

	*resolvedGeocodingToLoad = model.ResolvedGeocodingInput{
		District:     func(s string) *string { return &s }(""),
		Province:     func(s string) *string { return &s }(""),
		Municipality: func(s string) *string { return &s }(""),
		ZipCode:      func(s string) *string { return &s }(""),
		StreetName:   func(s string) *string { return &s }(""),
	}
	for _, addrComponent := range geocodeResults[0].AddressComponents {
		if funk.Contains(
			addrComponent.Types,
			"administrative_area_level_3",
		) {
			*resolvedGeocodingToLoad.District = addrComponent.LongName
		}

		if funk.Contains(
			addrComponent.Types,
			"administrative_area_level_2",
		) {
			*resolvedGeocodingToLoad.Municipality = addrComponent.LongName
		}

		if funk.Contains(
			addrComponent.Types,
			"administrative_area_level_1",
		) {
			*resolvedGeocodingToLoad.Province = addrComponent.LongName
		}

		if funk.Contains(
			addrComponent.Types,
			"postal_code",
		) {
			*resolvedGeocodingToLoad.ZipCode = addrComponent.LongName
		}

		if funk.Contains(
			addrComponent.Types,
			"street_address",
		) {
			*resolvedGeocodingToLoad.StreetName = addrComponent.LongName
		}
	}

	regionGroupKeyword := *resolvedGeocodingToLoad.Municipality
	splitMunicipality := strings.Split(*resolvedGeocodingToLoad.Municipality, " ")
	if splitMunicipality[0] == "Kota" || splitMunicipality[0] == "Kabupaten" {
		regionGroupKeyword = strings.Join(
			splitMunicipality[1:],
			" ",
		)
	}
	addressRegion, err := addrLoader.addressRegionDataSource.GetMongoDataSource().FindOne(
		map[string]interface{}{
			"cities": strings.ToUpper(regionGroupKeyword),
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
