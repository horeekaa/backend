package databasemoudatasourceinterfaces

import (
	mongodbmoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/mongodb/interfaces"
)

type MouDataSource interface {
	GetMongoDataSource() mongodbmoudatasourceinterfaces.MouDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbmoudatasourceinterfaces.MouDataSourceMongo) bool
}
