package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type CreateAccountFromAuthDataDependency struct{}

func (_ *CreateAccountFromAuthDataDependency) Bind() {
	container.Singleton(
		func(
			personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
		) accountdomainrepositoryinterfaces.CreateAccountFromAuthDataTransactionComponent {
			createAccFrmAuthDataComponent, _ :=
				accountdomainrepositories.NewCreateAccountFromAuthDataTransactionComponent(
					personDataSource,
					accountDataSource,
				)
			return createAccFrmAuthDataComponent
		},
	)

	container.Transient(
		func(
			usecaseComponent accountdomainrepositoryinterfaces.CreateAccountFromAuthDataTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) accountdomainrepositoryinterfaces.CreateAccountFromAuthDataRepository {
			createAccFromAuthDataRepo, _ := accountdomainrepositories.NewCreateAccountFromAuthDataRepository(
				usecaseComponent,
				mongoDBTransaction,
			)
			return createAccFromAuthDataRepo
		},
	)
}
