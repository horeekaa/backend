package googlemapcoreoperationinterfaces

import (
	"context"

	"googlemaps.github.io/maps"
)

type GoogleMapBasicOperation interface {
	ReverseGeocode(
		ctx context.Context, geocodingReq *maps.GeocodingRequest,
	) ([]maps.GeocodingResult, error)
}
