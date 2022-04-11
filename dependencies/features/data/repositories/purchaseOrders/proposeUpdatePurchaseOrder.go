package purchaseorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	databasepurchaseorderItemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	databasepurchaseOrderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositories "github.com/horeekaa/backend/features/purchaseOrders/data/repositories"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/utils"
)

type ProposeUpdatePurchaseOrderDependency struct{}

func (_ *ProposeUpdatePurchaseOrderDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderDataSource databasepurchaseOrderdatasourceinterfaces.PurchaseOrderDataSource,
			purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mouDataSource databasemoudatasourceinterfaces.MouDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
			purchaseOrderDataLoader purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader,
		) purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderTransactionComponent {
			proposeUpdatepurchaseOrderComponent, _ := purchaseorderdomainrepositories.NewProposeUpdatePurchaseOrderTransactionComponent(
				purchaseOrderDataSource,
				purchaseOrderItemDataSource,
				loggingDataSource,
				mouDataSource,
				mapProcessorUtility,
				purchaseOrderDataLoader,
			)
			return proposeUpdatepurchaseOrderComponent
		},
	)

	container.Transient(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			purchaseOrderDataSource databasepurchaseOrderdatasourceinterfaces.PurchaseOrderDataSource,
			purchaseOrderItemDataSource databasepurchaseorderItemdatasourceinterfaces.PurchaseOrderItemDataSource,
			trxComponent purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderTransactionComponent,
			createPurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent,
			proposeUpdatePurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent,
			approvePurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent,
			updateInvoiceTrxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
			createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderRepository {
			proposeUpdatepurchaseOrderRepo, _ := purchaseorderdomainrepositories.NewProposeUpdatePurchaseOrderRepository(
				memberAccessDataSource,
				purchaseOrderDataSource,
				purchaseOrderItemDataSource,
				trxComponent,
				createPurchaseOrderItemComponent,
				proposeUpdatePurchaseOrderItemComponent,
				approvePurchaseOrderItemComponent,
				updateInvoiceTrxComponent,
				createNotificationComponent,
				mongoDBTransaction,
			)
			return proposeUpdatepurchaseOrderRepo
		},
	)
}
