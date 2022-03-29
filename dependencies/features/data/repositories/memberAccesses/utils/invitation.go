package memberaccessdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositoryutilities "github.com/horeekaa/backend/features/memberAccesses/data/repositories/utils"
	memberaccessdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/utils"
)

type InvitationPayloadLoaderDependency struct{}

func (_ *InvitationPayloadLoaderDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
		) memberaccessdomainrepositoryutilityinterfaces.InvitationPayloadLoader {
			invitationPayloadLoader, _ := memberaccessdomainrepositoryutilities.NewInvitationPayloadLoader(
				accountDataSource,
				personDataSource,
			)
			return invitationPayloadLoader
		},
	)
}
