package moudomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	moudomainrepositoryutilities "github.com/horeekaa/backend/features/mous/data/repositories/utils"
	moudomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories/utils"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
)

type PartyLoaderDependency struct{}

func (_ *PartyLoaderDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
		) moudomainrepositoryutilityinterfaces.PartyLoader {
			partyLoader, _ := moudomainrepositoryutilities.NewPartyLoader(
				accountDataSource,
				personDataSource,
				organizationDataSource,
			)
			return partyLoader
		},
	)
}
