package databasesupplyorderdatasources

import (
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	mongodbsupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/mongodb/interfaces"
)

type supplyOrderDataSource struct {
	supplyOrderDataSourceRepoMongo mongodbsupplyorderdatasourceinterfaces.SupplyOrderDataSourceMongo
}

func (supplyOrderDataSource *supplyOrderDataSource) SetMongoDataSource(mongoDataSource mongodbsupplyorderdatasourceinterfaces.SupplyOrderDataSourceMongo) bool {
	supplyOrderDataSource.supplyOrderDataSourceRepoMongo = mongoDataSource
	return true
}

func (supplyOrderDataSource *supplyOrderDataSource) GetMongoDataSource() mongodbsupplyorderdatasourceinterfaces.SupplyOrderDataSourceMongo {
	return supplyOrderDataSource.supplyOrderDataSourceRepoMongo
}

func NewSupplyOrderDataSource() (databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource, error) {
	return &supplyOrderDataSource{}, nil
}
