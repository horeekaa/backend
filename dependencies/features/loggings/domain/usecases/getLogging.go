package loggingpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	loggingpresentationusecases "github.com/horeekaa/backend/features/loggings/domain/usecases"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
)

type GetLoggingUsecaseDependency struct{}

func (_ GetLoggingUsecaseDependency) bind() {
	container.Singleton(
		func(
			getLoggingRepo loggingdomainrepositoryinterfaces.GetLoggingRepository,
		) loggingpresentationusecaseinterfaces.GetLoggingUsecase {
			getLoggingUsecase, _ := loggingpresentationusecases.NewGetLoggingUsecase(
				getLoggingRepo,
			)
			return getLoggingUsecase
		},
	)
}
