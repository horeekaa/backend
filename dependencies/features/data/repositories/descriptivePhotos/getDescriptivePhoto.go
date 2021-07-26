package descriptivephotodomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	descriptivephotodomainrepositories "github.com/horeekaa/backend/features/descriptivePhotos/data/repositories"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
)

type GetDescriptivePhotoDependency struct{}

func (_ *GetDescriptivePhotoDependency) Bind() {
	container.Singleton(
		func(
			descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
		) descriptivephotodomainrepositoryinterfaces.GetDescriptivePhotoRepository {
			getDescriptivePhotoRepo, _ := descriptivephotodomainrepositories.NewGetDescriptivePhotoRepository(
				descriptivePhotoDataSource,
			)
			return getDescriptivePhotoRepo
		},
	)
}
