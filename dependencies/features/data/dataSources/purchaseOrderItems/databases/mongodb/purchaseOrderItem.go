package mongodbpurchaseorderitemdatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	mongodbpurchaseorderitemdatasources "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/mongodb"
	mongodbpurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/mongodb/interfaces"
	databasepurchaseorderitemdatasources "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/sources"
)

type PurchaseOrderItemDataSourceDependency struct{}

func (purcOrderItemDataSourceDpdcy *PurchaseOrderItemDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbpurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSourceMongo {
			purchaseOrderItemDataSourceMongo, _ := mongodbpurchaseorderitemdatasources.NewPurchaseOrderItemDataSourceMongo(basicOperation)
			return purchaseOrderItemDataSourceMongo
		},
	)

	container.Singleton(
		func(purchaseOrderItemDataSourceMongo mongodbpurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSourceMongo) databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource {
			purchaseOrderItemRepo, _ := databasepurchaseorderitemdatasources.NewPurchaseOrderItemDataSource()
			purchaseOrderItemRepo.SetMongoDataSource(purchaseOrderItemDataSourceMongo)
			return purchaseOrderItemRepo
		},
	)
}
