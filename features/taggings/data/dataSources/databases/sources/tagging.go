package databasetaggingdatasources

import (
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	mongodbtaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/mongodb/interfaces"
)

type taggingDataSource struct {
	taggingDataSourceRepoMongo mongodbtaggingdatasourceinterfaces.TaggingDataSourceMongo
}

func (taggingDataSource *taggingDataSource) SetMongoDataSource(mongoDataSource mongodbtaggingdatasourceinterfaces.TaggingDataSourceMongo) bool {
	taggingDataSource.taggingDataSourceRepoMongo = mongoDataSource
	return true
}

func (taggingDataSource *taggingDataSource) GetMongoDataSource() mongodbtaggingdatasourceinterfaces.TaggingDataSourceMongo {
	return taggingDataSource.taggingDataSourceRepoMongo
}

func NewTaggingDataSource() (databasetaggingdatasourceinterfaces.TaggingDataSource, error) {
	return &taggingDataSource{}, nil
}
