package databasedescriptivephotodatasources

import (
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	mongodbdescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/mongodb/interfaces"
)

type descriptivePhotoDataSource struct {
	descriptivePhotoDataSourceRepoMongo mongodbdescriptivephotodatasourceinterfaces.DescriptivePhotoDataSourceMongo
}

func (descPhotoDataSource *descriptivePhotoDataSource) SetMongoDataSource(mongoDataSource mongodbdescriptivephotodatasourceinterfaces.DescriptivePhotoDataSourceMongo) bool {
	descPhotoDataSource.descriptivePhotoDataSourceRepoMongo = mongoDataSource
	return true
}

func (descPhotoDataSource *descriptivePhotoDataSource) GetMongoDataSource() mongodbdescriptivephotodatasourceinterfaces.DescriptivePhotoDataSourceMongo {
	return descPhotoDataSource.descriptivePhotoDataSourceRepoMongo
}

func NewDescriptivePhotoDataSource() (databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource, error) {
	return &descriptivePhotoDataSource{}, nil
}
