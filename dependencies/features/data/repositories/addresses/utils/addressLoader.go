package addressdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	googlemapcoreoperationinterfaces "github.com/horeekaa/backend/core/maps/google/interfaces/operations"
	databaseaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	addressdomainrepositoryutilities "github.com/horeekaa/backend/features/addresses/data/repositories/utils"
	addressdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories/utils"
)

type AddressLoaderDependency struct{}

func (_ *AddressLoaderDependency) Bind() {
	container.Singleton(
		func(
			gMapOperation googlemapcoreoperationinterfaces.GoogleMapBasicOperation,
			addressRegionDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource,
		) addressdomainrepositoryutilityinterfaces.AddressLoader {
			addressLoader, _ := addressdomainrepositoryutilities.NewAddressLoader(
				gMapOperation,
				addressRegionDataSource,
			)
			return addressLoader
		},
	)
}
