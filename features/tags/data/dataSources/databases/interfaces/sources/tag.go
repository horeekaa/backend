package databasetagdatasourceinterfaces

import (
	mongodbtagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/mongodb/interfaces"
)

type TagDataSource interface {
	GetMongoDataSource() mongodbtagdatasourceinterfaces.TagDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbtagdatasourceinterfaces.TagDataSourceMongo) bool
}
