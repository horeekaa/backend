package purchaseorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositories "github.com/horeekaa/backend/features/purchaseOrderItems/data/repositories"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/utils"
)

type UpdatePurchaseOrderItemDependency struct{}

func (_ *UpdatePurchaseOrderItemDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
			purchaseOrderItemLoader purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader,
		) purchaseorderitemdomainrepositoryinterfaces.UpdatePurchaseOrderItemTransactionComponent {
			updatePurchaseOrderItemComponent, _ := purchaseorderitemdomainrepositories.NewUpdatePurchaseOrderItemTransactionComponent(
				purchaseOrderItemDataSource,
				purchaseOrderItemLoader,
			)
			return updatePurchaseOrderItemComponent
		},
	)
}