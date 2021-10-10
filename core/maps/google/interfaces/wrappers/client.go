package googlemapcorewrapperinterfaces

import (
	"context"

	"googlemaps.github.io/maps"
)

type GoogleMapClientWrapper interface {
	ReverseGeocode(
		ctx context.Context, geocodingReq *maps.GeocodingRequest,
	) ([]maps.GeocodingResult, error)
}
