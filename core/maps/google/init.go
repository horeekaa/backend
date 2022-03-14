package googlemapcoreclients

import (
	coreconfigs "github.com/horeekaa/backend/core/commons/configs"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	googlemapcoreclientinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/init"
	googlemapcorewrapperinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/wrappers"
	googlemapcorewrappers "github.com/horeekaa/backend/core/maps/google/wrappers"
	maps "googlemaps.github.io/maps"
)

type googleMapClient struct {
	client       googlemapcorewrapperinterfaces.GoogleMapClientWrapper
	pathIdentity string
}

func (gMapClient *googleMapClient) Initialize() (bool, error) {
	mapClient, err := maps.NewClient(
		maps.WithAPIKey(
			coreconfigs.GetEnvVariable(coreconfigs.GoogleAPIKey),
		),
	)
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			gMapClient.pathIdentity,
			err,
		)
	}

	wrappedClient, _ := googlemapcorewrappers.NewGoogleMapClientWrapper(mapClient)
	gMapClient.client = wrappedClient

	return true, nil
}

func (gMapClient *googleMapClient) GetGoogleMapClient() (googlemapcorewrapperinterfaces.GoogleMapClientWrapper, error) {
	if gMapClient.client == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			gMapClient.pathIdentity,
			nil,
		)
	}
	return gMapClient.client, nil
}

func NewGoogleMapClient() (googlemapcoreclientinterfaces.GoogleMapClient, error) {
	return &googleMapClient{
		pathIdentity: "GoogleMapClient",
	}, nil
}
