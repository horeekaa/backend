package memberaccessrefpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefpresentationusecases "github.com/horeekaa/backend/features/memberAccessRefs/domain/usecases"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
)

type CreateMemberAccessRefUsecaseDependency struct{}

func (_ *CreateMemberAccessRefUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
			createMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefRepository,
			logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
		) memberaccessrefpresentationusecaseinterfaces.CreateMemberAccessRefUsecase {
			memberAccRefUcase, _ := memberaccessrefpresentationusecases.NewCreateMemberAccessRefUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getPersonDataFromAccountRepo,
				createMemberAccessRefRepo,
				logEntityProposalActivityRepo,
			)
			return memberAccRefUcase
		},
	)
}
