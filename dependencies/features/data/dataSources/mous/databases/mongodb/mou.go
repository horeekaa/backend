package mongodbmoudatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	mongodbmoudatasources "github.com/horeekaa/backend/features/mous/data/dataSources/databases/mongodb"
	mongodbmoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/mongodb/interfaces"
	databasemoudatasources "github.com/horeekaa/backend/features/mous/data/dataSources/databases/sources"
)

type MouDataSourceDependency struct{}

func (orgDataSourceDpdcy *MouDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbmoudatasourceinterfaces.MouDataSourceMongo {
			mouDataSourceMongo, _ := mongodbmoudatasources.NewMouDataSourceMongo(basicOperation)
			return mouDataSourceMongo
		},
	)

	container.Singleton(
		func(mouDataSourceMongo mongodbmoudatasourceinterfaces.MouDataSourceMongo) databasemoudatasourceinterfaces.MouDataSource {
			mouRepo, _ := databasemoudatasources.NewMouDataSource()
			mouRepo.SetMongoDataSource(mouDataSourceMongo)
			return mouRepo
		},
	)
}
