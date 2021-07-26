package descriptivephotodomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	googlecloudstoragecoreoperationinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/operations"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	descriptivephotodomainrepositories "github.com/horeekaa/backend/features/descriptivePhotos/data/repositories"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
)

type UpdateDescriptivePhotoDependency struct{}

func (_ *UpdateDescriptivePhotoDependency) Bind() {
	container.Singleton(
		func(
			descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
			gcsBasicImageStoring googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation,
		) descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent {
			updateDescriptivePhotoComponent, _ := descriptivephotodomainrepositories.NewUpdateDescriptivePhotoTransactionComponent(
				descriptivePhotoDataSource,
				gcsBasicImageStoring,
			)
			return updateDescriptivePhotoComponent
		},
	)
}
