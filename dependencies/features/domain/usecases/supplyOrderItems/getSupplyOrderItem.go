package supplyorderitempresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitempresentationusecases "github.com/horeekaa/backend/features/supplyOrderItems/domain/usecases"
	supplyorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases"
)

type GetSupplyOrderItemUsecaseDependency struct{}

func (_ GetSupplyOrderItemUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getSupplyOrderItemRepo supplyorderitemdomainrepositoryinterfaces.GetSupplyOrderItemRepository,
		) supplyorderitempresentationusecaseinterfaces.GetSupplyOrderItemUsecase {
			getSupplyOrderItemUsecase, _ := supplyorderitempresentationusecases.NewGetSupplyOrderItemUsecase(
				getSupplyOrderItemRepo,
			)
			return getSupplyOrderItemUsecase
		},
	)
}
