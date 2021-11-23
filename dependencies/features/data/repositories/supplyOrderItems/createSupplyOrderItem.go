package supplyorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositories "github.com/horeekaa/backend/features/supplyOrderItems/data/repositories"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/utils"
)

type CreateSupplyOrderItemDependency struct{}

func (_ *CreateSupplyOrderItemDependency) Bind() {
	container.Singleton(
		func(
			supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			supplyOrderItemLoader supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader,
		) supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent {
			createSupplyOrderItemComponent, _ := supplyorderitemdomainrepositories.NewCreateSupplyOrderItemTransactionComponent(
				supplyOrderItemDataSource,
				loggingDataSource,
				supplyOrderItemLoader,
			)
			return createSupplyOrderItemComponent
		},
	)
}
