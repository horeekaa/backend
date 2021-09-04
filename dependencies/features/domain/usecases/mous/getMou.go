package moupresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moupresentationusecases "github.com/horeekaa/backend/features/mous/domain/usecases"
	moupresentationusecaseinterfaces "github.com/horeekaa/backend/features/mous/presentation/usecases"
)

type GetMouUsecaseDependency struct{}

func (_ GetMouUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getMouRepo moudomainrepositoryinterfaces.GetMouRepository,
		) moupresentationusecaseinterfaces.GetMouUsecase {
			getMouUsecase, _ := moupresentationusecases.NewGetMouUsecase(
				getMouRepo,
			)
			return getMouUsecase
		},
	)
}
