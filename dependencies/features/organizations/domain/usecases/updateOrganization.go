package organizationpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationpresentationusecases "github.com/horeekaa/backend/features/organizations/domain/usecases"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
)

type UpdateOrganizationUsecaseDependency struct{}

func (_ *UpdateOrganizationUsecaseDependency) Bind() {
	container.Singleton(
		func(
			manageAccountAuthenticationRepo accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository,
			getAccountMemberAccessRepo accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
			updateOrganizationRepo organizationdomainrepositoryinterfaces.UpdateOrganizationRepository,
			getOrganizationRepo organizationdomainrepositoryinterfaces.GetOrganizationRepository,
			logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
			logEntityApprovalActivityRepo loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository,
		) organizationpresentationusecaseinterfaces.UpdateOrganizationUsecase {
			updateOrganizationUsecase, _ := organizationpresentationusecases.NewUpdateOrganizationUsecase(
				manageAccountAuthenticationRepo,
				getAccountMemberAccessRepo,
				getPersonDataFromAccountRepo,
				updateOrganizationRepo,
				getOrganizationRepo,
				logEntityProposalActivityRepo,
				logEntityApprovalActivityRepo,
			)
			return updateOrganizationUsecase
		},
	)
}
