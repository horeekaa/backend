package notificationdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryutilities "github.com/horeekaa/backend/features/notifications/data/repositories/utils"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
)

type InvitationPayloadLoaderDependency struct{}

func (_ *InvitationPayloadLoaderDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
		) notificationdomainrepositoryutilityinterfaces.InvitationPayloadLoader {
			invitationPayloadLoader, _ := notificationdomainrepositoryutilities.NewInvitationPayloadLoader(
				accountDataSource,
				personDataSource,
			)
			return invitationPayloadLoader
		},
	)
}
