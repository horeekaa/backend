package descriptivephotopresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	descriptivephotopresentationusecases "github.com/horeekaa/backend/features/descriptivePhotos/domain/usecases"
	descriptivephotopresentationusecaseinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/presentation/usecases"
)

type GetDescriptivePhotoUsecaseDependency struct{}

func (_ GetDescriptivePhotoUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getDescriptivePhotoRepo descriptivephotodomainrepositoryinterfaces.GetDescriptivePhotoRepository,
		) descriptivephotopresentationusecaseinterfaces.GetDescriptivePhotoUsecase {
			getDescriptivePhotoUsecase, _ := descriptivephotopresentationusecases.NewGetDescriptivePhotoUsecase(
				getDescriptivePhotoRepo,
			)
			return getDescriptivePhotoUsecase
		},
	)
}
