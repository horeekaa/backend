package notificationdomainrepositoryutilityloaderdependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryloaderutility "github.com/horeekaa/backend/features/notifications/data/repositories/utils/payloadLoaders"
	notificationdomainrepositoryloaderutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils/payloadLoaders"
)

type InvitationPayloadLoaderDependency struct{}

func (_ *InvitationPayloadLoaderDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
		) notificationdomainrepositoryloaderutilityinterfaces.InvitationPayloadLoader {
			invitationPayloadLoader, _ := notificationdomainrepositoryloaderutility.NewInvitationPayloadLoader(
				accountDataSource,
				personDataSource,
			)
			return invitationPayloadLoader
		},
	)
}
