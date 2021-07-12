package memberaccessrefpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefpresentationusecases "github.com/horeekaa/backend/features/memberAccessRefs/domain/usecases"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type CreateMemberAccessRefUsecaseDependency struct{}

func (_ *CreateMemberAccessRefUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			createMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefRepository,
		) memberaccessrefpresentationusecaseinterfaces.CreateMemberAccessRefUsecase {
			memberAccessRefRefUcase, _ := memberaccessrefpresentationusecases.NewCreateMemberAccessRefUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				createMemberAccessRefRepo,
			)
			return memberAccessRefRefUcase
		},
	)
}
