package memberaccessdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositories "github.com/horeekaa/backend/features/memberAccesses/data/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type GetAllMemberAccessDependency struct{}

func (_ *GetAllMemberAccessDependency) Bind() {
	container.Singleton(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
		) memberaccessdomainrepositoryinterfaces.GetAllMemberAccessRepository {
			getAllMemberAccessRepo, _ := memberaccessdomainrepositories.NewGetAllMemberAccessRepository(
				memberAccessDataSource,
			)
			return getAllMemberAccessRepo
		},
	)
}
