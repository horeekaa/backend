package supplyorderitempresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitempresentationusecases "github.com/horeekaa/backend/features/supplyOrderItems/domain/usecases"
	supplyorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases"
)

type GetAllSupplyOrderItemUsecaseDependency struct{}

func (_ GetAllSupplyOrderItemUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllSupplyOrderItemRepo supplyorderitemdomainrepositoryinterfaces.GetAllSupplyOrderItemRepository,
		) supplyorderitempresentationusecaseinterfaces.GetAllSupplyOrderItemUsecase {
			getAllSupplyOrderItemUcase, _ := supplyorderitempresentationusecases.NewGetAllSupplyOrderItemUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getAllSupplyOrderItemRepo,
			)
			return getAllSupplyOrderItemUcase
		},
	)
}
