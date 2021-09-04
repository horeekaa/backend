package moupresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moupresentationusecases "github.com/horeekaa/backend/features/mous/domain/usecases"
	moupresentationusecaseinterfaces "github.com/horeekaa/backend/features/mous/presentation/usecases"
)

type CreateMouUsecaseDependency struct{}

func (_ *CreateMouUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			createMouRepo moudomainrepositoryinterfaces.CreateMouRepository,
		) moupresentationusecaseinterfaces.CreateMouUsecase {
			createMouUcase, _ := moupresentationusecases.NewCreateMouUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				createMouRepo,
			)
			return createMouUcase
		},
	)
}
