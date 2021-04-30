package organizationpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationpresentationusecases "github.com/horeekaa/backend/features/organizations/domain/usecases"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
)

type GetAllOrganizationUsecaseDependency struct{}

func (_ GetAllOrganizationUsecaseDependency) bind() {
	container.Singleton(
		func(
			manageAccountAuthenticationRepo accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository,
			getAccountOrganizationpo accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllOrganizationRepo organizationdomainrepositoryinterfaces.GetAllOrganizationRepository,
		) organizationpresentationusecaseinterfaces.GetAllOrganizationUsecase {
			getAllOrganizationUcase, _ := organizationpresentationusecases.NewGetAllOrganizationUsecase(
				manageAccountAuthenticationRepo,
				getAccountOrganizationpo,
				getAllOrganizationRepo,
			)
			return getAllOrganizationUcase
		},
	)
}
