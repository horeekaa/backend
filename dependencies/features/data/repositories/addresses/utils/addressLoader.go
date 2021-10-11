package addressdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	googlemapcoreoperationinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/operations"
	addressdomainrepositoryutilities "github.com/horeekaa/backend/features/addresses/data/repositories/utils"
	addressdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories/utils"
)

type AddressLoaderDependency struct{}

func (_ *AddressLoaderDependency) Bind() {
	container.Singleton(
		func(
			gMapOperation googlemapcoreoperationinterfaces.GoogleMapBasicOperation,
		) addressdomainrepositoryutilityinterfaces.AddressLoader {
			addressLoader, _ := addressdomainrepositoryutilities.NewAddressLoader(
				gMapOperation,
			)
			return addressLoader
		},
	)
}
