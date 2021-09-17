package tagdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositories "github.com/horeekaa/backend/features/tags/data/repositories"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
)

type GetAllTagDependency struct{}

func (_ *GetAllTagDependency) Bind() {
	container.Singleton(
		func(
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) tagdomainrepositoryinterfaces.GetAllTagRepository {
			getAllTagRepo, _ := tagdomainrepositories.NewGetAllTagRepository(
				tagDataSource,
				mongoQueryBuilder,
			)
			return getAllTagRepo
		},
	)
}
