package mouitempresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitempresentationusecases "github.com/horeekaa/backend/features/mouItems/domain/usecases"
	mouitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/mouItems/presentation/usecases"
)

type GetMouItemUsecaseDependency struct{}

func (_ GetMouItemUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getMouItemRepo mouitemdomainrepositoryinterfaces.GetMouItemRepository,
		) mouitempresentationusecaseinterfaces.GetMouItemUsecase {
			getMouItemUsecase, _ := mouitempresentationusecases.NewGetMouItemUsecase(
				getMouItemRepo,
			)
			return getMouItemUsecase
		},
	)
}
