package descriptivephotodomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	googlecloudstoragecoreoperationinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/operations"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	descriptivephotodomainrepositories "github.com/horeekaa/backend/features/descriptivePhotos/data/repositories"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
)

type CreateDescriptivePhotoDependency struct{}

func (_ *CreateDescriptivePhotoDependency) Bind() {
	container.Singleton(
		func(
			descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			gcsBasicImageStoring googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation,
		) descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent {
			createDescriptivePhotoComponent, _ := descriptivephotodomainrepositories.NewCreateDescriptivePhotoTransactionComponent(
				descriptivePhotoDataSource,
				loggingDataSource,
				gcsBasicImageStoring,
			)
			return createDescriptivePhotoComponent
		},
	)
}
