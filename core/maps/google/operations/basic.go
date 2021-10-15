package googlemapcoreoperation

import (
	"context"

	"googlemaps.github.io/maps"

	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	googlemapcoreclientinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/init"

	googlemapcoreoperationinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/operations"
)

type googleMapBasicOperation struct {
	client googlemapcoreclientinterfaces.GoogleMapClient
}

func (gMapBasicOperation *googleMapBasicOperation) ReverseGeocode(
	ctx context.Context, geocodingReq *maps.GeocodingRequest,
) ([]maps.GeocodingResult, error) {
	gMapClient, err := gMapBasicOperation.client.GetGoogleMapClient()
	if err != nil {
		return nil, err
	}

	results, err := gMapClient.ReverseGeocode(ctx, (*maps.GeocodingRequest)(geocodingReq))
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ReverseGeocodeFailed,
			"/googleMapBasicOperation/reverseGeocode",
			err,
		)
	}
	return results, nil
}

func NewGoogleMapBasicOperation(
	client googlemapcoreclientinterfaces.GoogleMapClient,
) (googlemapcoreoperationinterfaces.GoogleMapBasicOperation, error) {
	return &googleMapBasicOperation{
		client,
	}, nil
}
