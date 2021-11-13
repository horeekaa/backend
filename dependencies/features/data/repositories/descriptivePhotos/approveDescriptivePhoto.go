package descriptivephotodomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	googlecloudstoragecoreoperationinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/operations"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	descriptivephotodomainrepositories "github.com/horeekaa/backend/features/descriptivePhotos/data/repositories"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
)

type ApproveUpdateDescriptivePhotoDependency struct{}

func (_ *ApproveUpdateDescriptivePhotoDependency) Bind() {
	container.Singleton(
		func(
			descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent,
			gcsBasicImageStoring googlecloudstoragecoreoperationinterfaces.GCSBasicImageStoringOperation,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent {
			approveUpdateDescriptivePhotoComponent, _ := descriptivephotodomainrepositories.NewApproveUpdateDescriptivePhotoTransactionComponent(
				descriptivePhotoDataSource,
				loggingDataSource,
				updateDescriptivePhotoComponent,
				gcsBasicImageStoring,
				mapProcessorUtility,
			)
			return approveUpdateDescriptivePhotoComponent
		},
	)
}
