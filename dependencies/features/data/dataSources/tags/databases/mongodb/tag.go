package mongodbtagdatasourcedependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	mongodbtagdatasources "github.com/horeekaa/backend/features/tags/data/dataSources/databases/mongodb"
	mongodbtagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/mongodb/interfaces"
	databasetagdatasources "github.com/horeekaa/backend/features/tags/data/dataSources/databases/sources"
)

type TagDataSourceDependency struct{}

func (_ *TagDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbtagdatasourceinterfaces.TagDataSourceMongo {
			tagDataSourceMongo, _ := mongodbtagdatasources.NewTagDataSourceMongo(basicOperation)
			return tagDataSourceMongo
		},
	)

	container.Singleton(
		func(tagDataSourceMongo mongodbtagdatasourceinterfaces.TagDataSourceMongo) databasetagdatasourceinterfaces.TagDataSource {
			tagRepo, _ := databasetagdatasources.NewTagDataSource()
			tagRepo.SetMongoDataSource(tagDataSourceMongo)
			return tagRepo
		},
	)
}
