package mongodbnotificationdatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	mongodbnotificationdatasources "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/mongodb"
	mongodbnotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/mongodb/interfaces"
	databasenotificationdatasources "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/sources"
)

type NotificationDataSourceDependency struct{}

func (_ *NotificationDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbnotificationdatasourceinterfaces.NotificationDataSourceMongo {
			notificationDataSourceMongo, _ := mongodbnotificationdatasources.NewNotificationDataSourceMongo(basicOperation)
			return notificationDataSourceMongo
		},
	)

	container.Singleton(
		func(notificationDataSourceMongo mongodbnotificationdatasourceinterfaces.NotificationDataSourceMongo) databasenotificationdatasourceinterfaces.NotificationDataSource {
			notificationRepo, _ := databasenotificationdatasources.NewNotificationDataSource()
			notificationRepo.SetMongoDataSource(notificationDataSourceMongo)
			return notificationRepo
		},
	)
}
