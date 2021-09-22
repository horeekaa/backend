package purchaseorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositories "github.com/horeekaa/backend/features/purchaseOrderItems/data/repositories"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
)

type GetPurchaseOrderItemDependency struct{}

func (_ *GetPurchaseOrderItemDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
		) purchaseorderitemdomainrepositoryinterfaces.GetPurchaseOrderItemRepository {
			getPurchaseOrderItemRepo, _ := purchaseorderitemdomainrepositories.NewGetPurchaseOrderItemRepository(
				purchaseOrderItemDataSource,
			)
			return getPurchaseOrderItemRepo
		},
	)
}
