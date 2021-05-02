package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type GetAccountDependency struct{}

func (_ *GetAccountDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
		) accountdomainrepositoryinterfaces.GetAccountRepository {
			getAccountRepo, _ := accountdomainrepositories.NewGetAccountRepository(
				accountDataSource,
			)
			return getAccountRepo
		},
	)
}
