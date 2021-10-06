package purchaseorderpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderpresentationusecases "github.com/horeekaa/backend/features/purchaseOrders/domain/usecases"
	purchaseorderpresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrders/presentation/usecases"
)

type UpdatePurchaseOrderUsecaseDependency struct{}

func (_ *UpdatePurchaseOrderUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			proposeUpdatepurchaseOrderRepo purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderRepository,
			approveUpdatepurchaseOrderRepo purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderRepository,
		) purchaseorderpresentationusecaseinterfaces.UpdatePurchaseOrderUsecase {
			updatePurchaseOrderUsecase, _ := purchaseorderpresentationusecases.NewUpdatePurchaseOrderUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				proposeUpdatepurchaseOrderRepo,
				approveUpdatepurchaseOrderRepo,
			)
			return updatePurchaseOrderUsecase
		},
	)
}
