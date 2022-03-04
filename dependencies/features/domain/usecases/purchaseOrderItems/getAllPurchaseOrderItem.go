package purchaseorderitempresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitempresentationusecases "github.com/horeekaa/backend/features/purchaseOrderItems/domain/usecases"
	purchaseorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/presentation/usecases"
)

type GetAllPurchaseOrderItemUsecaseDependency struct{}

func (_ GetAllPurchaseOrderItemUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllPurchaseOrderItemRepo purchaseorderitemdomainrepositoryinterfaces.GetAllPurchaseOrderItemRepository,
		) purchaseorderitempresentationusecaseinterfaces.GetAllPurchaseOrderItemUsecase {
			getAllPurchaseOrderItemUcase, _ := purchaseorderitempresentationusecases.NewGetAllPurchaseOrderItemUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getAllPurchaseOrderItemRepo,
			)
			return getAllPurchaseOrderItemUcase
		},
	)
}
