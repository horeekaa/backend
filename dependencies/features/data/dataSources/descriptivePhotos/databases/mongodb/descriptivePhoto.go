package mongodbdescriptivephotodatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	mongodbdescriptivephotodatasources "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/mongodb"
	mongodbdescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/mongodb/interfaces"
	databasedescriptivephotodatasources "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/sources"
)

type DescriptivePhotoDataSourceDependency struct{}

func (_ *DescriptivePhotoDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbdescriptivephotodatasourceinterfaces.DescriptivePhotoDataSourceMongo {
			descriptivePhotoDataSourceMongo, _ := mongodbdescriptivephotodatasources.NewDescriptivePhotoDataSourceMongo(basicOperation)
			return descriptivePhotoDataSourceMongo
		},
	)

	container.Singleton(
		func(descriptivePhotoDataSourceMongo mongodbdescriptivephotodatasourceinterfaces.DescriptivePhotoDataSourceMongo) databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource {
			descriptivePhotoDataSource, _ := databasedescriptivephotodatasources.NewDescriptivePhotoDataSource()
			descriptivePhotoDataSource.SetMongoDataSource(descriptivePhotoDataSourceMongo)
			return descriptivePhotoDataSource
		},
	)
}
