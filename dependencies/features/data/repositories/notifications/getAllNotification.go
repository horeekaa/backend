package notificationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositories "github.com/horeekaa/backend/features/notifications/data/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	notificationdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories/utils"
)

type GetAllNotificationDependency struct{}

func (_ *GetAllNotificationDependency) Bind() {
	container.Singleton(
		func(
			notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
			notifLocalizationBuilder notificationdomainrepositoryutilityinterfaces.NotificationLocalizationBuilder,
		) notificationdomainrepositoryinterfaces.GetAllNotificationRepository {
			getAllNotificationRepo, _ := notificationdomainrepositories.NewGetAllNotificationRepository(
				notificationDataSource,
				notifLocalizationBuilder,
			)
			return getAllNotificationRepo
		},
	)
}
