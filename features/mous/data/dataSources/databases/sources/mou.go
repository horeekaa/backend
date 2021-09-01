package databasemoudatasources

import (
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	mongodbmoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/mongodb/interfaces"
)

type mouDataSource struct {
	mouDataSourceRepoMongo mongodbmoudatasourceinterfaces.MouDataSourceMongo
}

func (mouDataSource *mouDataSource) SetMongoDataSource(mongoDataSource mongodbmoudatasourceinterfaces.MouDataSourceMongo) bool {
	mouDataSource.mouDataSourceRepoMongo = mongoDataSource
	return true
}

func (mouDataSource *mouDataSource) GetMongoDataSource() mongodbmoudatasourceinterfaces.MouDataSourceMongo {
	return mouDataSource.mouDataSourceRepoMongo
}

func NewMouDataSource() (databasemoudatasourceinterfaces.MouDataSource, error) {
	return &mouDataSource{}, nil
}
