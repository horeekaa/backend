package mongodbsupplyorderdatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	mongodbsupplyorderdatasources "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/mongodb"
	mongodbsupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/mongodb/interfaces"
	databasesupplyorderdatasources "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/sources"
)

type SupplyOrderDataSourceDependency struct{}

func (_ *SupplyOrderDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbsupplyorderdatasourceinterfaces.SupplyOrderDataSourceMongo {
			supplyOrderDataSourceMongo, _ := mongodbsupplyorderdatasources.NewSupplyOrderDataSourceMongo(basicOperation)
			return supplyOrderDataSourceMongo
		},
	)

	container.Singleton(
		func(supplyOrderDataSourceMongo mongodbsupplyorderdatasourceinterfaces.SupplyOrderDataSourceMongo) databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource {
			supplyOrderRepo, _ := databasesupplyorderdatasources.NewSupplyOrderDataSource()
			supplyOrderRepo.SetMongoDataSource(supplyOrderDataSourceMongo)
			return supplyOrderRepo
		},
	)
}
