package purchaseordertosupplypresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	purchaseordertosupplypresentationusecases "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/usecases"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
)

type GetPurchaseOrderToSupplyUsecaseDependency struct{}

func (_ GetPurchaseOrderToSupplyUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getPurchaseOrderToSupplyRepo purchaseordertosupplydomainrepositoryinterfaces.GetPurchaseOrderToSupplyRepository,
		) purchaseordertosupplypresentationusecaseinterfaces.GetPurchaseOrderToSupplyUsecase {
			getPurchaseOrderToSupplyUsecase, _ := purchaseordertosupplypresentationusecases.NewGetPurchaseOrderToSupplyUsecase(
				getPurchaseOrderToSupplyRepo,
			)
			return getPurchaseOrderToSupplyUsecase
		},
	)
}
