package memberaccessrefpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefpresentationusecases "github.com/horeekaa/backend/features/memberAccessRefs/domain/usecases"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type UpdateMemberAccessRefUsecaseDependency struct{}

func (_ *UpdateMemberAccessRefUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			proposeUpdateMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefRepository,
			approveUpdateMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefRepository,
		) memberaccessrefpresentationusecaseinterfaces.UpdateMemberAccessRefUsecase {
			updateMemberAccessRefUsecase, _ := memberaccessrefpresentationusecases.NewUpdateMemberAccessRefUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				proposeUpdateMemberAccessRefRepo,
				approveUpdateMemberAccessRefRepo,
			)
			return updateMemberAccessRefUsecase
		},
	)
}
