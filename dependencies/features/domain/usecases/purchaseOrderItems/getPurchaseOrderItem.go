package purchaseorderItempresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitempresentationusecases "github.com/horeekaa/backend/features/purchaseOrderItems/domain/usecases"
	purchaseorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases"
)

type GetPurchaseOrderItemUsecaseDependency struct{}

func (_ GetPurchaseOrderItemUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getPurchaseOrderItemRepo purchaseorderitemdomainrepositoryinterfaces.GetPurchaseOrderItemRepository,
		) purchaseorderitempresentationusecaseinterfaces.GetPurchaseOrderItemUsecase {
			getPurchaseOrderItemUsecase, _ := purchaseorderitempresentationusecases.NewGetPurchaseOrderItemUsecase(
				getPurchaseOrderItemRepo,
			)
			return getPurchaseOrderItemUsecase
		},
	)
}
