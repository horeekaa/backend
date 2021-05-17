package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type GetAccountMemberAccessDependency struct{}

func (getAccountMemberAccessDependency *GetAccountMemberAccessDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			memberAccessDataSource databaseaccountdatasourceinterfaces.MemberAccessDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository {
			getAccountMemberAccessDataSource, _ := accountdomainrepositories.NewGetAccountMemberAccessRepository(
				accountDataSource,
				memberAccessDataSource,
				mapProcessorUtility,
			)
			return getAccountMemberAccessDataSource
		},
	)
}
