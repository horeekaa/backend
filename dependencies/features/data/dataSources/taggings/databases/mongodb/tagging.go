package mongodbtaggingdatasourcedependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	mongodbtaggingdatasources "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/mongodb"
	mongodbtaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/mongodb/interfaces"
	databasetaggingdatasources "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/sources"
)

type TaggingDataSourceDependency struct{}

func (_ *TaggingDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbtaggingdatasourceinterfaces.TaggingDataSourceMongo {
			taggingDataSourceMongo, _ := mongodbtaggingdatasources.NewTaggingDataSourceMongo(basicOperation)
			return taggingDataSourceMongo
		},
	)

	container.Singleton(
		func(taggingDataSourceMongo mongodbtaggingdatasourceinterfaces.TaggingDataSourceMongo) databasetaggingdatasourceinterfaces.TaggingDataSource {
			taggingRepo, _ := databasetaggingdatasources.NewTaggingDataSource()
			taggingRepo.SetMongoDataSource(taggingDataSourceMongo)
			return taggingRepo
		},
	)
}
