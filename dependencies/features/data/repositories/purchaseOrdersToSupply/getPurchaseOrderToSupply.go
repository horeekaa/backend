package purchaseordertosupplydomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasepurchaseordertoSupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	purchaseordertosupplydomainrepositories "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/repositories"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
)

type GetPurchaseOrderToSupplyDependency struct{}

func (_ *GetPurchaseOrderToSupplyDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderToSupplyDataSource databasepurchaseordertoSupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
		) purchaseordertosupplydomainrepositoryinterfaces.GetPurchaseOrderToSupplyRepository {
			getPurchaseOrderToSupplyRepo, _ := purchaseordertosupplydomainrepositories.NewGetPurchaseOrderToSupplyRepository(
				purchaseOrderToSupplyDataSource,
			)
			return getPurchaseOrderToSupplyRepo
		},
	)
}
