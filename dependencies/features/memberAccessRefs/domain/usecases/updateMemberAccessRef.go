package memberaccessrefpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefpresentationusecases "github.com/horeekaa/backend/features/memberAccessRefs/domain/usecases"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
)

type UpdateMemberAccessRefUsecaseDependency struct{}

func (_ *UpdateMemberAccessRefUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getPersonDataFromAccountRepo accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
			updateMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefRepository,
			getMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.GetMemberAccessRefRepository,
			logEntityProposalActivityRepo loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository,
			logEntityApprovalActivityRepo loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository,
		) memberaccessrefpresentationusecaseinterfaces.UpdateMemberAccessRefUsecase {
			updateMemberAccessRefUsecase, _ := memberaccessrefpresentationusecases.NewUpdateMemberAccessRefUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getPersonDataFromAccountRepo,
				updateMemberAccessRefRepo,
				getMemberAccessRefRepo,
				logEntityProposalActivityRepo,
				logEntityApprovalActivityRepo,
			)
			return updateMemberAccessRefUsecase
		},
	)
}
