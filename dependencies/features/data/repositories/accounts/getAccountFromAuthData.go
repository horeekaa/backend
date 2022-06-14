package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	firebaseauthdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/authentication/interfaces"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
)

type GetAccountFromAuthDataDependency struct{}

func (_ *GetAccountFromAuthDataDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			firebaseDataSource firebaseauthdatasourceinterfaces.FirebaseAuthRepo,
		) accountdomainrepositoryinterfaces.GetAccountFromAuthData {
			getAccFromAuthDataRepo, _ := accountdomainrepositories.NewGetAccountFromAuthDataRepository(
				accountDataSource,
				memberAccessDataSource,
				firebaseDataSource,
			)
			return getAccFromAuthDataRepo
		},
	)
}
