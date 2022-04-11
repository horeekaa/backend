package purchaseorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
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
			mouDataSource databasemoudatasourceinterfaces.MouDataSource,
			purchaseOrderDataLoader purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader,
		) purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderTransactionComponent {
			createPurchaseOrderComponent, _ := purchaseorderdomainrepositories.NewCreatePurchaseOrderTransactionComponent(
				purchaseOrderDataSource,
				loggingDataSource,
				mouDataSource,
				purchaseOrderDataLoader,
			)
			return createPurchaseOrderComponent
		},
	)

	container.Transient(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			trxComponent purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderTransactionComponent,
			createPurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent,
			createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderRepository {
			createPurchaseOrderRepo, _ := purchaseorderdomainrepositories.NewCreatePurchaseOrderRepository(
				memberAccessDataSource,
				trxComponent,
				createPurchaseOrderItemComponent,
				createNotificationComponent,
				mongoDBTransaction,
			)
			return createPurchaseOrderRepo
		},
	)
}
