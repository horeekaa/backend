package moupresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moupresentationusecases "github.com/horeekaa/backend/features/mous/domain/usecases"
	moupresentationusecaseinterfaces "github.com/horeekaa/backend/features/mous/presentation/usecases"
)

type UpdateMouUsecaseDependency struct{}

func (_ *UpdateMouUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			proposeUpdateMouRepo moudomainrepositoryinterfaces.ProposeUpdateMouRepository,
			approveUpdateMouRepo moudomainrepositoryinterfaces.ApproveUpdateMouRepository,
		) moupresentationusecaseinterfaces.UpdateMouUsecase {
			updateMouUsecase, _ := moupresentationusecases.NewUpdateMouUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				proposeUpdateMouRepo,
				approveUpdateMouRepo,
			)
			return updateMouUsecase
		},
	)
}
