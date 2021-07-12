package organizationpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationpresentationusecases "github.com/horeekaa/backend/features/organizations/domain/usecases"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
)

type UpdateOrganizationUsecaseDependency struct{}

func (_ *UpdateOrganizationUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			proposeUpdateOrganizationRepo organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationRepository,
			approveUpdateOrganizationRepo organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationRepository,
			createMemberAccessRepo memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository,
			getOrganizationRepo organizationdomainrepositoryinterfaces.GetOrganizationRepository,
		) organizationpresentationusecaseinterfaces.UpdateOrganizationUsecase {
			updateOrganizationUsecase, _ := organizationpresentationusecases.NewUpdateOrganizationUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				proposeUpdateOrganizationRepo,
				approveUpdateOrganizationRepo,
				createMemberAccessRepo,
				getOrganizationRepo,
			)
			return updateOrganizationUsecase
		},
	)
}
