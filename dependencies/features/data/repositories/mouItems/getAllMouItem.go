package mouitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositories "github.com/horeekaa/backend/features/mouItems/data/repositories"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
)

type GetAllMouItemDependency struct{}

func (_ *GetAllMouItemDependency) Bind() {
	container.Singleton(
		func(
			mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) mouitemdomainrepositoryinterfaces.GetAllMouItemRepository {
			getAllMouItemRepo, _ := mouitemdomainrepositories.NewGetAllMouItemRepository(
				mouItemDataSource,
				mongoQueryBuilder,
			)
			return getAllMouItemRepo
		},
	)
}
