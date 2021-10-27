package notificationdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	notificationdomainrepositoryutilities "github.com/horeekaa/backend/features/notifications/data/repositories/utils"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
	notificationdomainrepositoryloaderutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils/payloadLoaders"
)

type MasterPayloadLoaderDependency struct{}

func (_ *MasterPayloadLoaderDependency) Bind() {
	container.Singleton(
		func(
			invitationPayloadLoader notificationdomainrepositoryloaderutilityinterfaces.InvitationPayloadLoader,
		) notificationdomainrepositoryutilityinterfaces.MasterPayloadLoader {
			masterPayloadLoader, _ := notificationdomainrepositoryutilities.NewMasterPayloadLoader(
				invitationPayloadLoader,
			)
			return masterPayloadLoader
		},
	)
}
