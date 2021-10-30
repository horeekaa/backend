package purchaseordertosupplypresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	purchaseordertosupplypresentationusecases "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/usecases"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
)

type ProcessPurchaseOrderToSupplyUsecaseDependency struct{}

func (_ *ProcessPurchaseOrderToSupplyUsecaseDependency) Bind() {
	container.Singleton(
		func(
			processPurchaseOrderToSupplyRepo purchaseordertosupplydomainrepositoryinterfaces.ProcessPurchaseOrderToSupplyRepository,
		) purchaseordertosupplypresentationusecaseinterfaces.ProcessPurchaseOrderToSupplyUsecase {
			purchaseOrderToSupplyRefUcase, _ := purchaseordertosupplypresentationusecases.NewProcessPurchaseOrderToSupplyUsecase(
				processPurchaseOrderToSupplyRepo,
			)
			return purchaseOrderToSupplyRefUcase
		},
	)
}
