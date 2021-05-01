package organizationpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationpresentationusecases "github.com/horeekaa/backend/features/organizations/domain/usecases"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
)

type GetOrganizationUsecaseDependency struct{}

func (_ GetOrganizationUsecaseDependency) bind() {
	container.Singleton(
		func(
			getOrganizationRepo organizationdomainrepositoryinterfaces.GetOrganizationRepository,
		) organizationpresentationusecaseinterfaces.GetOrganizationUsecase {
			getOrganizationUsecase, _ := organizationpresentationusecases.NewGetOrganizationUsecase(
				getOrganizationRepo,
			)
			return getOrganizationUsecase
		},
	)
}
