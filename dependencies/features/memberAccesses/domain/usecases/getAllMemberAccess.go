package memberaccesspresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccesspresentationusecases "github.com/horeekaa/backend/features/memberAccesses/domain/usecases"
	memberaccesspresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases"
)

type GetAllMemberAccessUsecaseDependency struct{}

func (_ GetAllMemberAccessUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAllMemberAccessRepository,
		) memberaccesspresentationusecaseinterfaces.GetAllMemberAccessUsecase {
			getAllMemberAccessUcase, _ := memberaccesspresentationusecases.NewGetAllMemberAccessUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getAllMemberAccessRepo,
			)
			return getAllMemberAccessUcase
		},
	)
}
