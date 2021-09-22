package databasepurchaseorderItemdatasourceinterfaces

import (
	mongodbpurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/mongodb/interfaces"
)

type PurchaseOrderItemDataSource interface {
	GetMongoDataSource() mongodbpurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbpurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSourceMongo) bool
}