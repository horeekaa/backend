package purchaseordertosupplypresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	purchaseordertosupplypresentationusecases "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/usecases"
	purchaseordertosupplypresentationusecaseinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/presentation/usecases"
)

type GetAllPurchaseOrderToSupplyUsecaseDependency struct{}

func (_ GetAllPurchaseOrderToSupplyUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllPurchaseOrderToSupplyRepo purchaseordertosupplydomainrepositoryinterfaces.GetAllPurchaseOrderToSupplyRepository,
		) purchaseordertosupplypresentationusecaseinterfaces.GetAllPurchaseOrderToSupplyUsecase {
			getAllPurchaseOrderToSupplyUcase, _ := purchaseordertosupplypresentationusecases.NewGetAllPurchaseOrderToSupplyUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getAllPurchaseOrderToSupplyRepo,
			)
			return getAllPurchaseOrderToSupplyUcase
		},
	)
}
