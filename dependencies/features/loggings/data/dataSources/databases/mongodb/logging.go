package mongodbloggingdatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	mongodbloggingdatasources "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/mongodb"
	mongodbloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/mongodb/interfaces"
	databaseloggingdatasources "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/sources"
)

type LoggingDataSourceDependency struct{}

func (logDataSourceDpdcy *LoggingDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbloggingdatasourceinterfaces.LoggingDataSourceMongo {
			loggingDataSourceMongo, _ := mongodbloggingdatasources.NewLoggingDataSourceMongo(basicOperation)
			return loggingDataSourceMongo
		},
	)

	container.Singleton(
		func(loggingDataSourceMongo mongodbloggingdatasourceinterfaces.LoggingDataSourceMongo) databaseloggingdatasourceinterfaces.LoggingDataSource {
			loggingDataSource, _ := databaseloggingdatasources.NewLoggingDataSource()
			loggingDataSource.SetMongoDataSource(loggingDataSourceMongo)
			return loggingDataSource
		},
	)
}
