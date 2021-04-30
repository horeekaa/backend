package memberaccessrefdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	memberaccessrefdomainrepositories "github.com/horeekaa/backend/features/memberAccessRefs/data/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
)

type GetMemberAccessRefDependency struct{}

func (_ *GetMemberAccessRefDependency) bind() {
	container.Singleton(
		func(
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
		) memberaccessrefdomainrepositoryinterfaces.GetMemberAccessRefRepository {
			getMemberAccessRefRepo, _ := memberaccessrefdomainrepositories.NewGetMemberAccessRefRepository(
				memberAccessRefDataSource,
			)
			return getMemberAccessRefRepo
		},
	)
}
