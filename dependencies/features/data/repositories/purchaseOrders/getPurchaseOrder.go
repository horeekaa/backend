package purchaseorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasepurchaseOrderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositories "github.com/horeekaa/backend/features/purchaseOrders/data/repositories"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
)

type GetPurchaseOrderDependency struct{}

func (_ *GetPurchaseOrderDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderDataSource databasepurchaseOrderdatasourceinterfaces.PurchaseOrderDataSource,
		) purchaseorderdomainrepositoryinterfaces.GetPurchaseOrderRepository {
			getPurchaseOrderRepo, _ := purchaseorderdomainrepositories.NewGetPurchaseOrderRepository(
				purchaseOrderDataSource,
			)
			return getPurchaseOrderRepo
		},
	)
}
