package databasenotificationdatasources

import (
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	mongodbnotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/mongodb/interfaces"
)

type notificationDataSource struct {
	notificationDataSourceRepoMongo mongodbnotificationdatasourceinterfaces.NotificationDataSourceMongo
}

func (notificationDataSource *notificationDataSource) SetMongoDataSource(mongoDataSource mongodbnotificationdatasourceinterfaces.NotificationDataSourceMongo) bool {
	notificationDataSource.notificationDataSourceRepoMongo = mongoDataSource
	return true
}

func (notificationDataSource *notificationDataSource) GetMongoDataSource() mongodbnotificationdatasourceinterfaces.NotificationDataSourceMongo {
	return notificationDataSource.notificationDataSourceRepoMongo
}

func NewNotificationDataSource() (databasenotificationdatasourceinterfaces.NotificationDataSource, error) {
	return &notificationDataSource{}, nil
}
