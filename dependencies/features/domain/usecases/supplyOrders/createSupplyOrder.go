package supplyorderpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderpresentationusecases "github.com/horeekaa/backend/features/supplyOrders/domain/usecases"
	supplyorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases"
)

type CreateSupplyOrderUsecaseDependency struct{}

func (_ *CreateSupplyOrderUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			createSupplyOrderRepo supplyorderdomainrepositoryinterfaces.CreateSupplyOrderRepository,
		) supplyorderpresentationusecaseinterfaces.CreateSupplyOrderUsecase {
			supplyOrderRefUcase, _ := supplyorderpresentationusecases.NewCreateSupplyOrderUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				createSupplyOrderRepo,
			)
			return supplyOrderRefUcase
		},
	)
}
