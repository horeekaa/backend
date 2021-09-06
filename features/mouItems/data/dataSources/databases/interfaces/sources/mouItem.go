package databasemouitemdatasourceinterfaces

import (
	mongodbmouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/mongodb/interfaces"
)

type MouItemDataSource interface {
	GetMongoDataSource() mongodbmouitemdatasourceinterfaces.MouItemDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbmouitemdatasourceinterfaces.MouItemDataSourceMongo) bool
}
