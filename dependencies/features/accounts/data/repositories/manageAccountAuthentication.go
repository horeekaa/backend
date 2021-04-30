package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	firebaseauthdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/authentication/interfaces"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type ManageAccountAuthenticationDependency struct{}

func (manageAccAuthDependency *ManageAccountAuthenticationDependency) bind() {
	container.Singleton(
		func(
			personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			firebaseDataSource firebaseauthdatasourceinterfaces.FirebaseAuthRepo,
		) accountdomainrepositoryinterfaces.ManageAccountAuthenticationTransactionComponent {
			manageAccAuthTransactionComponent, _ :=
				accountdomainrepositories.NewManageAccountAuthenticationTransactionComponent(
					personDataSource,
					accountDataSource,
					firebaseDataSource,
				)
			return manageAccAuthTransactionComponent
		},
	)

	container.Transient(
		func(
			usecaseComponent accountdomainrepositoryinterfaces.ManageAccountAuthenticationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository {
			accountAuthenticationRepo, _ := accountdomainrepositories.NewManageAccountAuthenticationRepository(
				usecaseComponent,
				mongoDBTransaction,
			)
			return accountAuthenticationRepo
		},
	)
}
