package memberaccessrefpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefpresentationusecases "github.com/horeekaa/backend/features/memberAccessRefs/domain/usecases"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
)

type GetAllMemberAccessRefUsecaseDependency struct{}

func (_ GetAllMemberAccessRefUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.GetAllMemberAccessRefRepository,
		) memberaccessrefpresentationusecaseinterfaces.GetAllMemberAccessRefUsecase {
			getAllMemberAccessRefUcase, _ := memberaccessrefpresentationusecases.NewGetAllMemberAccessRefUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getAllMemberAccessRefRepo,
			)
			return getAllMemberAccessRefUcase
		},
	)
}
