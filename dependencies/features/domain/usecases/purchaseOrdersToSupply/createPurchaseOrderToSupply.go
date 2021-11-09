package purchaseordertosupplypresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	purchaseordertosupplypresentationusecases "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/usecases"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
)

type CreatePurchaseOrderToSupplyUsecaseDependency struct{}

func (_ *CreatePurchaseOrderToSupplyUsecaseDependency) Bind() {
	container.Singleton(
		func(
			createPurchaseOrderToSupplyRepo purchaseordertosupplydomainrepositoryinterfaces.CreatePurchaseOrderToSupplyRepository,
		) purchaseordertosupplypresentationusecaseinterfaces.CreatePurchaseOrderToSupplyUsecase {
			purchaseOrderToSupplyRefUcase, _ := purchaseordertosupplypresentationusecases.NewCreatePurchaseOrderToSupplyUsecase(
				createPurchaseOrderToSupplyRepo,
			)
			return purchaseOrderToSupplyRefUcase
		},
	)
}
