package purchaseorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositories "github.com/horeekaa/backend/features/purchaseOrderItems/data/repositories"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/utils"
)

type ProposeUpdatePurchaseOrderItemDependency struct{}

func (_ *ProposeUpdatePurchaseOrderItemDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			purchaseOrderItemLoader purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent {
			proposeUpdatePurchaseOrderItemComponent, _ := purchaseorderitemdomainrepositories.NewProposeUpdatePurchaseOrderItemTransactionComponent(
				purchaseOrderItemDataSource,
				loggingDataSource,
				purchaseOrderItemLoader,
				mapProcessorUtility,
			)
			return proposeUpdatePurchaseOrderItemComponent
		},
	)
}
