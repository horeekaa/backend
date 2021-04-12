package databaseloggingdatasources

import (
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	mongodbloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/mongodb/interfaces"
)

type loggingDataSource struct {
	loggingDataSourceRepoMongo mongodbloggingdatasourceinterfaces.LoggingDataSourceMongo
}

func (loggingDataSource *loggingDataSource) SetMongoDataSource(mongoDataSource mongodbloggingdatasourceinterfaces.LoggingDataSourceMongo) bool {
	loggingDataSource.loggingDataSourceRepoMongo = mongoDataSource
	return true
}

func (loggingDataSource *loggingDataSource) GetMongoDataSource() mongodbloggingdatasourceinterfaces.LoggingDataSourceMongo {
	return loggingDataSource.loggingDataSourceRepoMongo
}

func NewLoggingDataSource() (databaseloggingdatasourceinterfaces.LoggingDataSource, error) {
	return &loggingDataSource{}, nil
}
