package memberaccessdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositories "github.com/horeekaa/backend/features/memberAccesses/data/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type GetAllMemberAccessDependency struct{}

func (_ *GetAllMemberAccessDependency) Bind() {
	container.Singleton(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) memberaccessdomainrepositoryinterfaces.GetAllMemberAccessRepository {
			getAllMemberAccessRepo, _ := memberaccessdomainrepositories.NewGetAllMemberAccessRepository(
				memberAccessDataSource,
				mongoQueryBuilder,
			)
			return getAllMemberAccessRepo
		},
	)
}
