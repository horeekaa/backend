package databaseloggingdatasourceinterfaces

import (
	mongodbloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/mongodb/interfaces"
)

type LoggingDataSource interface {
	GetMongoDataSource() mongodbloggingdatasourceinterfaces.LoggingDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbloggingdatasourceinterfaces.LoggingDataSourceMongo) bool
}
