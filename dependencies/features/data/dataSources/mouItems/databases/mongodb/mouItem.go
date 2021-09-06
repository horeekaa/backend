package mongodbmouItemdatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mongodbmouitemdatasources "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/mongodb"
	mongodbmouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/mongodb/interfaces"
	databasemouitemdatasources "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/sources"
)

type MouItemDataSourceDependency struct{}

func (orgDataSourceDpdcy *MouItemDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbmouitemdatasourceinterfaces.MouItemDataSourceMongo {
			mouItemDataSourceMongo, _ := mongodbmouitemdatasources.NewMouItemDataSourceMongo(basicOperation)
			return mouItemDataSourceMongo
		},
	)

	container.Singleton(
		func(mouItemDataSourceMongo mongodbmouitemdatasourceinterfaces.MouItemDataSourceMongo) databasemouitemdatasourceinterfaces.MouItemDataSource {
			mouItemRepo, _ := databasemouitemdatasources.NewMouItemDataSource()
			mouItemRepo.SetMongoDataSource(mouItemDataSourceMongo)
			return mouItemRepo
		},
	)
}
