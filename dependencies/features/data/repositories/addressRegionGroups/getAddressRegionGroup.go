package addressregiongroupdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	addressregiongroupdomainrepositories "github.com/horeekaa/backend/features/addressRegionGroups/data/repositories"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
)

type GetAddressRegionGroupDependency struct{}

func (_ *GetAddressRegionGroupDependency) Bind() {
	container.Singleton(
		func(
			addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource,
		) addressregiongroupdomainrepositoryinterfaces.GetAddressRegionGroupRepository {
			getAddressRegionGroupRepo, _ := addressregiongroupdomainrepositories.NewGetAddressRegionGroupRepository(
				addressRegionGroupDataSource,
			)
			return getAddressRegionGroupRepo
		},
	)
}
