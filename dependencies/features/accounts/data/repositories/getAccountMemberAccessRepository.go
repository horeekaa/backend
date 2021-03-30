package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type GetAccountMemberAccessDependency struct{}

func (getAccountMemberAccessDependency *GetAccountMemberAccessDependency) bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			memberAccessDataSource databaseaccountdatasourceinterfaces.MemberAccessDataSource,
		) accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository {
			getAccountMemberAccessDataSource, _ := accountdomainrepositories.NewGetAccountMemberAccessRepository(
				accountDataSource,
				memberAccessDataSource,
			)
			return getAccountMemberAccessDataSource
		},
	)
}
