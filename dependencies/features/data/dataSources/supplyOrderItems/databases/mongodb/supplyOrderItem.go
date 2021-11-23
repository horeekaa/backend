package mongodbsupplyorderitemdatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	mongodbsupplyorderitemdatasources "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/mongodb"
	mongodbsupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/mongodb/interfaces"
	databasesupplyorderitemdatasources "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/sources"
)

type SupplyOrderItemDataSourceDependency struct{}

func (_ *SupplyOrderItemDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbsupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSourceMongo {
			supplyOrderItemDataSourceMongo, _ := mongodbsupplyorderitemdatasources.NewSupplyOrderItemDataSourceMongo(basicOperation)
			return supplyOrderItemDataSourceMongo
		},
	)

	container.Singleton(
		func(supplyOrderItemDataSourceMongo mongodbsupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSourceMongo) databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource {
			supplyOrderItemRepo, _ := databasesupplyorderitemdatasources.NewSupplyOrderItemDataSource()
			supplyOrderItemRepo.SetMongoDataSource(supplyOrderItemDataSourceMongo)
			return supplyOrderItemRepo
		},
	)
}
