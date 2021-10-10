package googlemapcoreclientinterfaces

import googlemapcorewrapperinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/wrappers"

type GoogleMapClient interface {
	Initialize() (bool, error)
	GetGoogleMapClient() (googlemapcorewrapperinterfaces.GoogleMapClientWrapper, error)
}
