package purchaseorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositories "github.com/horeekaa/backend/features/purchaseOrders/data/repositories"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/utils"
)

type CreatePurchaseOrderDependency struct{}

func (_ *CreatePurchaseOrderDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			purchaseOrderDataLoader purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader,
		) purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderTransactionComponent {
			createPurchaseOrderComponent, _ := purchaseorderdomainrepositories.NewCreatePurchaseOrderTransactionComponent(
				purchaseOrderDataSource,
				loggingDataSource,
				purchaseOrderDataLoader,
			)
			return createPurchaseOrderComponent
		},
	)

	container.Transient(
		func(
			trxComponent purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderTransactionComponent,
			createPurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderRepository {
			createPurchaseOrderRepo, _ := purchaseorderdomainrepositories.NewCreatePurchaseOrderRepository(
				trxComponent,
				createPurchaseOrderItemComponent,
				mongoDBTransaction,
			)
			return createPurchaseOrderRepo
		},
	)
}
