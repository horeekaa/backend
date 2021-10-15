package mongodbaddressdatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	mongodbaddressdatasources "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/mongodb"
	mongodbaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/mongodb/interfaces"
	databaseaddressdatasources "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/sources"
)

type AddressDataSourceDependency struct{}

func (_ *AddressDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbaddressdatasourceinterfaces.AddressDataSourceMongo {
			addressDataSourceMongo, _ := mongodbaddressdatasources.NewAddressDataSourceMongo(basicOperation)
			return addressDataSourceMongo
		},
	)

	container.Singleton(
		func(addressDataSourceMongo mongodbaddressdatasourceinterfaces.AddressDataSourceMongo) databaseaddressdatasourceinterfaces.AddressDataSource {
			addressDataSource, _ := databaseaddressdatasources.NewAddressDataSource()
			addressDataSource.SetMongoDataSource(addressDataSourceMongo)
			return addressDataSource
		},
	)
}
