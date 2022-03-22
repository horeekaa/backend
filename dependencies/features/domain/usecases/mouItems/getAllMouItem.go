package mouitempresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitempresentationusecases "github.com/horeekaa/backend/features/mouItems/domain/usecases"
	mouitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/mouItems/presentation/usecases"
)

type GetAllMouItemUsecaseDependency struct{}

func (_ GetAllMouItemUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllMouItemRepo mouitemdomainrepositoryinterfaces.GetAllMouItemRepository,
		) mouitempresentationusecaseinterfaces.GetAllMouItemUsecase {
			getAllMouItemUcase, _ := mouitempresentationusecases.NewGetAllMouItemUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getAllMouItemRepo,
			)
			return getAllMouItemUcase
		},
	)
}
