package databasenotificationdatasourceinterfaces

import (
	mongodbnotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/mongodb/interfaces"
)

type NotificationDataSource interface {
	GetMongoDataSource() mongodbnotificationdatasourceinterfaces.NotificationDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbnotificationdatasourceinterfaces.NotificationDataSourceMongo) bool
}
