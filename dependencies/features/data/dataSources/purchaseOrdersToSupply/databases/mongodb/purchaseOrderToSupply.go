package mongodbpurchaseordertosupplydatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	mongodbpurchaseordertosupplydatasources "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/mongodb"
	mongodbpurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/mongodb/interfaces"
	databasepurchaseordertosupplydatasources "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/sources"
)

type PurchaseOrderToSupplyDataSourceDependency struct{}

func (_ *PurchaseOrderToSupplyDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbpurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSourceMongo {
			purchaseOrderToSupplyDataSourceMongo, _ := mongodbpurchaseordertosupplydatasources.NewPurchaseOrderToSupplyDataSourceMongo(basicOperation)
			return purchaseOrderToSupplyDataSourceMongo
		},
	)

	container.Singleton(
		func(purchaseOrderToSupplyDataSourceMongo mongodbpurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSourceMongo) databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource {
			purchaseOrderToSupplyRepo, _ := databasepurchaseordertosupplydatasources.NewPurchaseOrderToSupplyDataSource()
			purchaseOrderToSupplyRepo.SetMongoDataSource(purchaseOrderToSupplyDataSourceMongo)
			return purchaseOrderToSupplyRepo
		},
	)
}
