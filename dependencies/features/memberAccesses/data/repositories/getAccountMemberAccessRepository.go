package memberaccessdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositories "github.com/horeekaa/backend/features/memberAccesses/data/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type GetAccountMemberAccessDependency struct{}

func (getAccountMemberAccessDependency *GetAccountMemberAccessDependency) Bind() {
	container.Singleton(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository {
			getAccountMemberAccessDataSource, _ := memberaccessdomainrepositories.NewGetAccountMemberAccessRepository(
				memberAccessDataSource,
				mapProcessorUtility,
			)
			return getAccountMemberAccessDataSource
		},
	)
}
