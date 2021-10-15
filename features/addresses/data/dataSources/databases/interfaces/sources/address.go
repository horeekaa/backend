package databaseaddressdatasourceinterfaces

import (
	mongodbaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/mongodb/interfaces"
)

type AddressDataSource interface {
	GetMongoDataSource() mongodbaddressdatasourceinterfaces.AddressDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbaddressdatasourceinterfaces.AddressDataSourceMongo) bool
}
