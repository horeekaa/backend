package supplyorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositories "github.com/horeekaa/backend/features/supplyOrders/data/repositories"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/utils"
)

type CreateSupplyOrderDependency struct{}

func (_ *CreateSupplyOrderDependency) Bind() {
	container.Singleton(
		func(
			supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			supplyOrderDataLoader supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader,
		) supplyorderdomainrepositoryinterfaces.CreateSupplyOrderTransactionComponent {
			createsupplyOrderComponent, _ := supplyorderdomainrepositories.NewCreateSupplyOrderTransactionComponent(
				supplyOrderDataSource,
				loggingDataSource,
				supplyOrderDataLoader,
			)
			return createsupplyOrderComponent
		},
	)

	container.Transient(
		func(
			trxComponent supplyorderdomainrepositoryinterfaces.CreateSupplyOrderTransactionComponent,
			createSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) supplyorderdomainrepositoryinterfaces.CreateSupplyOrderRepository {
			createsupplyOrderRepo, _ := supplyorderdomainrepositories.NewCreateSupplyOrderRepository(
				trxComponent,
				createSupplyOrderItemComponent,
				mongoDBTransaction,
			)
			return createsupplyOrderRepo
		},
	)
}
