package notificationdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryutilities "github.com/horeekaa/backend/features/notifications/data/repositories/utils"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
	notificationdomainrepositoryloaderutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils/payloadLoaders"
)

type MasterPayloadLoaderDependency struct{}

func (_ *MasterPayloadLoaderDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
			invitationPayloadLoader notificationdomainrepositoryloaderutilityinterfaces.InvitationPayloadLoader,
		) notificationdomainrepositoryutilityinterfaces.MasterPayloadLoader {
			masterPayloadLoader, _ := notificationdomainrepositoryutilities.NewMasterPayloadLoader(
				accountDataSource,
				personDataSource,
				invitationPayloadLoader,
			)
			return masterPayloadLoader
		},
	)
}
