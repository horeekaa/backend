package organizationpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationpresentationusecases "github.com/horeekaa/backend/features/organizations/domain/usecases"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
)

type CreateOrganizationUsecaseDependency struct{}

func (_ *CreateOrganizationUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
			createOrganizationRepo organizationdomainrepositoryinterfaces.CreateOrganizationRepository,
			logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
		) organizationpresentationusecaseinterfaces.CreateOrganizationUsecase {
			organizationRefUcase, _ := organizationpresentationusecases.NewCreateOrganizationUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getPersonDataFromAccountRepo,
				createOrganizationRepo,
				logEntityProposalActivityRepo,
			)
			return organizationRefUcase
		},
	)
}
