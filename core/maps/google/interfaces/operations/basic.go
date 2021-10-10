package googlemapcoreoperationinterfaces

import (
	"context"

	googlemapcoretypes "github.com/horeekaa/backend/core/maps/google/types"
)

type GoogleMapBasicOperation interface {
	ReverseGeocode(
		ctx context.Context, geocodingReq *googlemapcoretypes.GeocodingRequest,
	) (googlemapcoretypes.GeocodingResult, error)
}
