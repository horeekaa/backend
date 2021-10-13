package addressregiongroupdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databaseaddressRegionGroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	addressregiongroupdomainrepositories "github.com/horeekaa/backend/features/addressRegionGroups/data/repositories"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
)

type GetAllAddressRegionGroupDependency struct{}

func (_ *GetAllAddressRegionGroupDependency) Bind() {
	container.Singleton(
		func(
			addressRegionGroupDataSource databaseaddressRegionGroupdatasourceinterfaces.AddressRegionGroupDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) addressregiongroupdomainrepositoryinterfaces.GetAllAddressRegionGroupRepository {
			getAllAddressRegionGroupRepo, _ := addressregiongroupdomainrepositories.NewGetAllAddressRegionGroupRepository(
				addressRegionGroupDataSource,
				mongoQueryBuilder,
			)
			return getAllAddressRegionGroupRepo
		},
	)
}
