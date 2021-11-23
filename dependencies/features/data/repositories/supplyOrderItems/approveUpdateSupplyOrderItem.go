package supplyorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositories "github.com/horeekaa/backend/features/supplyOrderItems/data/repositories"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
)

type ApproveUpdateSupplyOrderItemDependency struct{}

func (_ *ApproveUpdateSupplyOrderItemDependency) Bind() {
	container.Singleton(
		func(
			supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
			purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent {
			approveUpdateSupplyOrderItemComponent, _ := supplyorderitemdomainrepositories.NewApproveUpdateSupplyOrderItemTransactionComponent(
				supplyOrderItemDataSource,
				purchaseOrderToSupplyDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdateSupplyOrderItemComponent
		},
	)
}
