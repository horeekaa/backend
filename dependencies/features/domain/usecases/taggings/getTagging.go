package taggingpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingpresentationusecases "github.com/horeekaa/backend/features/taggings/domain/usecases"
	taggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/taggings/presentation/usecases"
)

type GetTaggingUsecaseDependency struct{}

func (_ GetTaggingUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getTaggingRepo taggingdomainrepositoryinterfaces.GetTaggingRepository,
		) taggingpresentationusecaseinterfaces.GetTaggingUsecase {
			getTaggingUsecase, _ := taggingpresentationusecases.NewGetTaggingUsecase(
				getTaggingRepo,
			)
			return getTaggingUsecase
		},
	)
}
