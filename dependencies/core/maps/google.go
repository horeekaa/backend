package mapcoredependencies

import (
	"github.com/golobby/container/v2"
	googlemapcoreclients "github.com/horeekaa/backend/core/maps/google"
	googlemapcoreclientinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/init"
	googlemapcoreoperationinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/operations"
	googlemapcoreoperations "github.com/horeekaa/backend/core/maps/google/operations"
)

type GoogleMapDependency struct{}

func (_ GoogleMapDependency) Bind() {
	container.Singleton(
		func() googlemapcoreclientinterfaces.GoogleMapClient {
			gMapClient, _ := googlemapcoreclients.NewGoogleMapClient()
			gMapClient.Initialize()
			return gMapClient
		},
	)

	container.Singleton(
		func(
			gMapClient googlemapcoreclientinterfaces.GoogleMapClient,
		) googlemapcoreoperationinterfaces.GoogleMapBasicOperation {
			gMapBasicOps, _ := googlemapcoreoperations.NewGoogleMapBasicOperation(
				gMapClient,
			)
			return gMapBasicOps
		},
	)
}
