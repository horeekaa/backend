package databasepurchaseorderdatasourceinterfaces

import (
	mongodbpurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/mongodb/interfaces"
)

type PurchaseOrderDataSource interface {
	GetMongoDataSource() mongodbpurchaseorderdatasourceinterfaces.PurchaseOrderDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbpurchaseorderdatasourceinterfaces.PurchaseOrderDataSourceMongo) bool
}
