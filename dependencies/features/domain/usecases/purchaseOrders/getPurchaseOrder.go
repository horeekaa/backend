package purchaseorderpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderpresentationusecases "github.com/horeekaa/backend/features/purchaseOrders/domain/usecases"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
)

type GetPurchaseOrderUsecaseDependency struct{}

func (_ GetPurchaseOrderUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getPurchaseOrderRepo purchaseorderdomainrepositoryinterfaces.GetPurchaseOrderRepository,
		) purchaseorderpresentationusecaseinterfaces.GetPurchaseOrderUsecase {
			getPurchaseOrderUsecase, _ := purchaseorderpresentationusecases.NewGetPurchaseOrderUsecase(
				getPurchaseOrderRepo,
			)
			return getPurchaseOrderUsecase
		},
	)
}
