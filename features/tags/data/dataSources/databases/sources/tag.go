package databasetagdatasources

import (
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	mongodbtagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/mongodb/interfaces"
)

type tagDataSource struct {
	tagDataSourceRepoMongo mongodbtagdatasourceinterfaces.TagDataSourceMongo
}

func (tagDataSource *tagDataSource) SetMongoDataSource(mongoDataSource mongodbtagdatasourceinterfaces.TagDataSourceMongo) bool {
	tagDataSource.tagDataSourceRepoMongo = mongoDataSource
	return true
}

func (tagDataSource *tagDataSource) GetMongoDataSource() mongodbtagdatasourceinterfaces.TagDataSourceMongo {
	return tagDataSource.tagDataSourceRepoMongo
}

func NewTagDataSource() (databasetagdatasourceinterfaces.TagDataSource, error) {
	return &tagDataSource{}, nil
}
