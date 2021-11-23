package supplyorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositories "github.com/horeekaa/backend/features/supplyOrderItems/data/repositories"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/utils"
)

type ProposeUpdateSupplyOrderItemDependency struct{}

func (_ *ProposeUpdateSupplyOrderItemDependency) Bind() {
	container.Singleton(
		func(
			supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			supplyOrderItemLoader supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemTransactionComponent {
			proposeUpdateSupplyOrderItemComponent, _ := supplyorderitemdomainrepositories.NewProposeUpdateSupplyOrderItemTransactionComponent(
				supplyOrderItemDataSource,
				loggingDataSource,
				supplyOrderItemLoader,
				mapProcessorUtility,
			)
			return proposeUpdateSupplyOrderItemComponent
		},
	)
}
