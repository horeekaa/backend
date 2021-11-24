package databasesupplyorderdatasourceinterfaces

import (
	mongodbsupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/mongodb/interfaces"
)

type SupplyOrderDataSource interface {
	GetMongoDataSource() mongodbsupplyorderdatasourceinterfaces.SupplyOrderDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbsupplyorderdatasourceinterfaces.SupplyOrderDataSourceMongo) bool
}
