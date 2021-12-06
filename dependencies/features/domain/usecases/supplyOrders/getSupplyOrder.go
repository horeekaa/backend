package supplyorderpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderpresentationusecases "github.com/horeekaa/backend/features/supplyOrders/domain/usecases"
	supplyorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrders/presentation/usecases"
)

type GetSupplyOrderUsecaseDependency struct{}

func (_ GetSupplyOrderUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getSupplyOrderRepo supplyorderdomainrepositoryinterfaces.GetSupplyOrderRepository,
		) supplyorderpresentationusecaseinterfaces.GetSupplyOrderUsecase {
			getSupplyOrderUsecase, _ := supplyorderpresentationusecases.NewGetSupplyOrderUsecase(
				getSupplyOrderRepo,
			)
			return getSupplyOrderUsecase
		},
	)
}
