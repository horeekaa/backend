package googlemapcorewrappers

import (
	"context"

	googlemapcorewrapperinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/wrappers"
	maps "googlemaps.github.io/maps"
)

type googleMapClientWrapper struct {
	*maps.Client
}

func (gMap *googleMapClientWrapper) ReverseGeocode(
	ctx context.Context, geocodingReq *maps.GeocodingRequest,
) ([]maps.GeocodingResult, error) {
	result, err := gMap.Client.ReverseGeocode(ctx, geocodingReq)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NewGoogleMapClientWrapper(mapClientToWrap *maps.Client) (googlemapcorewrapperinterfaces.GoogleMapClientWrapper, error) {
	return &googleMapClientWrapper{
		mapClientToWrap,
	}, nil
}
