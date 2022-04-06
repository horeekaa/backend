package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type UpdateAccountDependency struct{}

func (_ *UpdateAccountDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
		) accountdomainrepositoryinterfaces.UpdateAccountTransactionComponent {
			updateAccountComponent, _ := accountdomainrepositories.NewUpdateAccountTransactionComponent(
				accountDataSource,
				personDataSource,
			)
			return updateAccountComponent
		},
	)

	container.Transient(
		func(
			updateAccountComponent accountdomainrepositoryinterfaces.UpdateAccountTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) accountdomainrepositoryinterfaces.UpdateAccountRepository {
			updateAccountRepo, _ := accountdomainrepositories.NewUpdateAccountRepository(
				updateAccountComponent,
				mongoDBTransaction,
			)
			return updateAccountRepo
		},
	)
}
