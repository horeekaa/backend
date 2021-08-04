package databasetaggingdatasourceinterfaces

import (
	mongodbtaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/mongodb/interfaces"
)

type TaggingDataSource interface {
	GetMongoDataSource() mongodbtaggingdatasourceinterfaces.TaggingDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbtaggingdatasourceinterfaces.TaggingDataSourceMongo) bool
}
