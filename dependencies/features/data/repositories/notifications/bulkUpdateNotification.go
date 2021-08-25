package notificationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositories "github.com/horeekaa/backend/features/notifications/data/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
)

type BulkUpdateNotificationDependency struct{}

func (_ *BulkUpdateNotificationDependency) Bind() {
	container.Singleton(
		func(
			notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
			notifLocalizationBuilder notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder,
		) notificationdomainrepositoryinterfaces.BulkUpdateNotificationTransactionComponent {
			bulkUpdateNotificationComponent, _ := notificationdomainrepositories.NewBulkUpdateNotificationTransactionComponent(
				notificationDataSource,
				notifLocalizationBuilder,
			)
			return bulkUpdateNotificationComponent
		},
	)

	container.Transient(
		func(
			bulkUpdateNotificationComponent notificationdomainrepositoryinterfaces.BulkUpdateNotificationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) notificationdomainrepositoryinterfaces.BulkUpdateNotificationRepository {
			bulkUpdateNotificationRepo, _ := notificationdomainrepositories.NewBulkUpdateNotificationRepository(
				bulkUpdateNotificationComponent,
				mongoDBTransaction,
			)
			return bulkUpdateNotificationRepo
		},
	)
}
