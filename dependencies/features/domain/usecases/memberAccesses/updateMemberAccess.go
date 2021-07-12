package memberaccesspresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccesspresentationusecases "github.com/horeekaa/backend/features/memberAccesses/domain/usecases"
	memberaccesspresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases"
)

type UpdateMemberAccessUsecaseDependency struct{}

func (_ *UpdateMemberAccessUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			proposeUpdateMemberAccessRepo memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessRepository,
			approveUpdateMemberAccessRepo memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessRepository,
		) memberaccesspresentationusecaseinterfaces.UpdateMemberAccessUsecase {
			updateMemberAccessUsecase, _ := memberaccesspresentationusecases.NewUpdateMemberAccessUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				proposeUpdateMemberAccessRepo,
				approveUpdateMemberAccessRepo,
			)
			return updateMemberAccessUsecase
		},
	)
}
