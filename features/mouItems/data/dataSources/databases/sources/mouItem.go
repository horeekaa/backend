package databasemouitemdatasources

import (
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mongodbmouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/mongodb/interfaces"
)

type mouItemDataSource struct {
	mouItemDataSourceRepoMongo mongodbmouitemdatasourceinterfaces.MouItemDataSourceMongo
}

func (mouItemDataSource *mouItemDataSource) SetMongoDataSource(mongoDataSource mongodbmouitemdatasourceinterfaces.MouItemDataSourceMongo) bool {
	mouItemDataSource.mouItemDataSourceRepoMongo = mongoDataSource
	return true
}

func (mouItemDataSource *mouItemDataSource) GetMongoDataSource() mongodbmouitemdatasourceinterfaces.MouItemDataSourceMongo {
	return mouItemDataSource.mouItemDataSourceRepoMongo
}

func NewMouItemDataSource() (databasemouitemdatasourceinterfaces.MouItemDataSource, error) {
	return &mouItemDataSource{}, nil
}
