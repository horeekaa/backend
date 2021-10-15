package addresspresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	addresspresentationusecases "github.com/horeekaa/backend/features/addresses/domain/usecases"
	addresspresentationusecaseinterfaces "github.com/horeekaa/backend/features/addresses/presentation/usecases"
)

type GetAddressUsecaseDependency struct{}

func (_ GetAddressUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAddressRepo addressdomainrepositoryinterfaces.GetAddressRepository,
		) addresspresentationusecaseinterfaces.GetAddressUsecase {
			getAddressUsecase, _ := addresspresentationusecases.NewGetAddressUsecase(
				getAddressRepo,
			)
			return getAddressUsecase
		},
	)
}
