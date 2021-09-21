package mongodbpurchaseorderdatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	mongodbpurchaseorderdatasources "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/mongodb"
	mongodbpurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/mongodb/interfaces"
	databasepurchaseorderdatasources "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/sources"
)

type PurchaseOrderDataSourceDependency struct{}

func (purcOrderDataSourceDpdcy *PurchaseOrderDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbpurchaseorderdatasourceinterfaces.PurchaseOrderDataSourceMongo {
			purchaseOrderDataSourceMongo, _ := mongodbpurchaseorderdatasources.NewPurchaseOrderDataSourceMongo(basicOperation)
			return purchaseOrderDataSourceMongo
		},
	)

	container.Singleton(
		func(purchaseOrderDataSourceMongo mongodbpurchaseorderdatasourceinterfaces.PurchaseOrderDataSourceMongo) databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource {
			purchaseOrderRepo, _ := databasepurchaseorderdatasources.NewPurchaseOrderDataSource()
			purchaseOrderRepo.SetMongoDataSource(purchaseOrderDataSourceMongo)
			return purchaseOrderRepo
		},
	)
}
