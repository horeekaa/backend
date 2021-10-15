package mongodbaddressregiongroupdatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databaseaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	mongodbaddressregiongroupdatasources "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/mongodb"
	mongodbaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/mongodb/interfaces"
	databaseaddressregiongroupdatasources "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/sources"
)

type AddressRegionGroupDataSourceDependency struct{}

func (_ *AddressRegionGroupDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSourceMongo {
			addressRegionGroupDataSourceMongo, _ := mongodbaddressregiongroupdatasources.NewAddressRegionGroupDataSourceMongo(basicOperation)
			return addressRegionGroupDataSourceMongo
		},
	)

	container.Singleton(
		func(addressRegionGroupDataSourceMongo mongodbaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSourceMongo) databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource {
			addressRegionGroupDataSource, _ := databaseaddressregiongroupdatasources.NewAddressRegionGroupDataSource()
			addressRegionGroupDataSource.SetMongoDataSource(addressRegionGroupDataSourceMongo)
			return addressRegionGroupDataSource
		},
	)
}
