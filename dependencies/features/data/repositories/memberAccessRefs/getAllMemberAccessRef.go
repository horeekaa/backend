package memberaccessrefdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositories "github.com/horeekaa/backend/features/memberAccessRefs/data/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
)

type GetAllMemberAccessRefDependency struct{}

func (_ *GetAllMemberAccessRefDependency) Bind() {
	container.Singleton(
		func(
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) memberaccessrefdomainrepositoryinterfaces.GetAllMemberAccessRefRepository {
			getAllMemberAccessRefRepo, _ := memberaccessrefdomainrepositories.NewGetAllMemberAccessRefRepository(
				memberAccessRefDataSource,
				mongoQueryBuilder,
			)
			return getAllMemberAccessRefRepo
		},
	)
}
