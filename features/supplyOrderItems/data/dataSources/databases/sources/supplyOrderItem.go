package databasesupplyorderitemdatasources

import (
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	mongodbsupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/mongodb/interfaces"
)

type supplyOrderItemDataSource struct {
	supplyOrderItemDataSourceRepoMongo mongodbsupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSourceMongo
}

func (supplyOrderItemDataSource *supplyOrderItemDataSource) SetMongoDataSource(mongoDataSource mongodbsupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSourceMongo) bool {
	supplyOrderItemDataSource.supplyOrderItemDataSourceRepoMongo = mongoDataSource
	return true
}

func (supplyOrderItemDataSource *supplyOrderItemDataSource) GetMongoDataSource() mongodbsupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSourceMongo {
	return supplyOrderItemDataSource.supplyOrderItemDataSourceRepoMongo
}

func NewSupplyOrderItemDataSource() (databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource, error) {
	return &supplyOrderItemDataSource{}, nil
}
