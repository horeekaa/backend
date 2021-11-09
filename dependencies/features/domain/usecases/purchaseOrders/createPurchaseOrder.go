package purchaseorderpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderpresentationusecases "github.com/horeekaa/backend/features/purchaseOrders/domain/usecases"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
)

type CreatePurchaseOrderUsecaseDependency struct{}

func (_ *CreatePurchaseOrderUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			createPurchaseOrderRepo purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderRepository,
		) purchaseorderpresentationusecaseinterfaces.CreatePurchaseOrderUsecase {
			purchaseOrderRefUcase, _ := purchaseorderpresentationusecases.NewCreatePurchaseOrderUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				createPurchaseOrderRepo,
			)
			return purchaseOrderRefUcase
		},
	)
}
