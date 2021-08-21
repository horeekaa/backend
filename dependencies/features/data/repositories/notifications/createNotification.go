package notificationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	firebasemessagingcoreoperationinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/operations"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositories "github.com/horeekaa/backend/features/notifications/data/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
)

type CreateNotificationDependency struct{}

func (_ *CreateNotificationDependency) Bind() {
	container.Singleton(
		func(
			notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
			firebaseMessaging firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation,
			notifLocalizationBuilder notificationdomainrepositoryinterfaces.NotificationLocalizationBuilder,
		) notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent {
			createNotificationComponent, _ := notificationdomainrepositories.NewCreateNotificationTransactionComponent(
				notificationDataSource,
				firebaseMessaging,
				notifLocalizationBuilder,
			)
			return createNotificationComponent
		},
	)
}
