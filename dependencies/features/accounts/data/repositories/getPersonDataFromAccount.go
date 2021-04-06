package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type GetPersonDataFromAccountDependency struct{}

func (getPersonDataFromAccountDependency *GetPersonDataFromAccountDependency) bind() {
	container.Singleton(
		func(
			personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
		) accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository {
			getPersonDataFromAccountRepo, _ := accountdomainrepositories.NewGetPersonDataFromAccountRepository(
				personDataSource,
				accountDataSource,
			)
			return getPersonDataFromAccountRepo
		},
	)
}
