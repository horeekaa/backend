package databasepurchaseorderdatasources

import (
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	mongodbpurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/mongodb/interfaces"
)

type purchaseOrderDataSource struct {
	purchaseOrderDataSourceRepoMongo mongodbpurchaseorderdatasourceinterfaces.PurchaseOrderDataSourceMongo
}

func (purchaseOrderDataSource *purchaseOrderDataSource) SetMongoDataSource(mongoDataSource mongodbpurchaseorderdatasourceinterfaces.PurchaseOrderDataSourceMongo) bool {
	purchaseOrderDataSource.purchaseOrderDataSourceRepoMongo = mongoDataSource
	return true
}

func (purchaseOrderDataSource *purchaseOrderDataSource) GetMongoDataSource() mongodbpurchaseorderdatasourceinterfaces.PurchaseOrderDataSourceMongo {
	return purchaseOrderDataSource.purchaseOrderDataSourceRepoMongo
}

func NewPurchaseOrderDataSource() (databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource, error) {
	return &purchaseOrderDataSource{}, nil
}
