package supplyorderpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderpresentationusecases "github.com/horeekaa/backend/features/supplyOrders/domain/usecases"
	supplyorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases"
)

type GetAllSupplyOrderUsecaseDependency struct{}

func (_ GetAllSupplyOrderUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllSupplyOrderRepo supplyorderdomainrepositoryinterfaces.GetAllSupplyOrderRepository,
		) supplyorderpresentationusecaseinterfaces.GetAllSupplyOrderUsecase {
			getAllSupplyOrderUcase, _ := supplyorderpresentationusecases.NewGetAllSupplyOrderUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getAllSupplyOrderRepo,
			)
			return getAllSupplyOrderUcase
		},
	)
}
