package notificationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	firebasemessagingcoreoperationinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/operations"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositories "github.com/horeekaa/backend/features/notifications/data/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
)

type CreateNotificationDependency struct{}

func (_ *CreateNotificationDependency) Bind() {
	container.Singleton(
		func(
			notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
			firebaseMessaging firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation,
			notifLocalizationBuilder notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder,
			masterPayloadLoader notificationdomainrepositoryutilityinterfaces.MasterPayloadLoader,
		) notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent {
			createNotificationComponent, _ := notificationdomainrepositories.NewCreateNotificationTransactionComponent(
				notificationDataSource,
				firebaseMessaging,
				notifLocalizationBuilder,
				masterPayloadLoader,
			)
			return createNotificationComponent
		},
	)
}
