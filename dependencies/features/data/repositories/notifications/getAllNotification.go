package notificationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositories "github.com/horeekaa/backend/features/notifications/data/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
)

type GetAllNotificationDependency struct{}

func (_ *GetAllNotificationDependency) Bind() {
	container.Singleton(
		func(
			notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
			notifLocalizationBuilder notificationdomainrepositoryinterfaces.NotificationLocalizationBuilder,
		) notificationdomainrepositoryinterfaces.GetAllNotificationRepository {
			getAllNotificationRepo, _ := notificationdomainrepositories.NewGetAllNotificationRepository(
				notificationDataSource,
				notifLocalizationBuilder,
			)
			return getAllNotificationRepo
		},
	)
}
