package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	authenticationcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/interfaces"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type GetAccountFromAuthDataDependency struct{}

func (_ *GetAccountFromAuthDataDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			firebaseDataSource authenticationcoreclientinterfaces.AuthenticationRepo,
		) accountdomainrepositoryinterfaces.GetAccountFromAuthData {
			getAccFromAuthDataRepo, _ := accountdomainrepositories.NewGetAccountFromAuthDataRepository(
				accountDataSource,
				firebaseDataSource,
			)
			return getAccFromAuthDataRepo
		},
	)
}
