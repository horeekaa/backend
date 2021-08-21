package notificationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositories "github.com/horeekaa/backend/features/notifications/data/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
)

type BulkUpdateNotificationDependency struct{}

func (_ *BulkUpdateNotificationDependency) Bind() {
	container.Singleton(
		func(
			notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
			notifLocalizationBuilder notificationdomainrepositoryinterfaces.NotificationLocalizationBuilder,
		) notificationdomainrepositoryinterfaces.BulkUpdateNotificationTransactionComponent {
			bulkUpdateNotificationComponent, _ := notificationdomainrepositories.NewBulkUpdateNotificationTransactionComponent(
				notificationDataSource,
				notifLocalizationBuilder,
			)
			return bulkUpdateNotificationComponent
		},
	)
}
