package tagpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	tagpresentationusecases "github.com/horeekaa/backend/features/tags/domain/usecases"
	tagpresentationusecaseinterfaces "github.com/horeekaa/backend/features/tags/presentation/usecases"
)

type GetTagUsecaseDependency struct{}

func (_ GetTagUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getTagRepo tagdomainrepositoryinterfaces.GetTagRepository,
		) tagpresentationusecaseinterfaces.GetTagUsecase {
			getTagUsecase, _ := tagpresentationusecases.NewGetTagUsecase(
				getTagRepo,
			)
			return getTagUsecase
		},
	)
}
