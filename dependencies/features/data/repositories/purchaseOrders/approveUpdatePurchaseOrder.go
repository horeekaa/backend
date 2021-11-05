package purchaseorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	databasepurchaseOrderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositories "github.com/horeekaa/backend/features/purchaseOrders/data/repositories"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/utils"
)

type ApproveUpdatePurchaseOrderDependency struct{}

func (_ *ApproveUpdatePurchaseOrderDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderDataSource databasepurchaseOrderdatasourceinterfaces.PurchaseOrderDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
			mouDataSource databasemoudatasourceinterfaces.MouDataSource,
			purchaseOrderDataLoader purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader,
		) purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderTransactionComponent {
			approveUpdatePurchaseOrderComponent, _ := purchaseorderdomainrepositories.NewApproveUpdatePurchaseOrderTransactionComponent(
				purchaseOrderDataSource,
				loggingDataSource,
				mapProcessorUtility,
				mouDataSource,
				purchaseOrderDataLoader,
			)
			return approveUpdatePurchaseOrderComponent
		},
	)

	container.Transient(
		func(
			trxComponent purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderRepository {
			approveUpdatePurchaseOrderRepo, _ := purchaseorderdomainrepositories.NewApproveUpdatePurchaseOrderRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return approveUpdatePurchaseOrderRepo
		},
	)
}
