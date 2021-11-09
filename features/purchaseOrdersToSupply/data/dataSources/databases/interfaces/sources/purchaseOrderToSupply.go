package databasepurchaseordertosupplydatasourceinterfaces

import (
	mongodbpurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/mongodb/interfaces"
)

type PurchaseOrderToSupplyDataSource interface {
	GetMongoDataSource() mongodbpurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbpurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSourceMongo) bool
}
