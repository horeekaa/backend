package purchaseordertosupplydomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	purchaseordertosupplydomainrepositories "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/repositories"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
)

type ProcessPurchaseOrderToSupplyDependency struct{}

func (_ *ProcessPurchaseOrderToSupplyDependency) Bind() {
	container.Singleton(
		func(
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
		) purchaseordertosupplydomainrepositoryinterfaces.ProcessPurchaseOrderToSupplyTransactionComponent {
			proposeUpdatepurchaseOrderToSupplyComponent, _ := purchaseordertosupplydomainrepositories.NewProcessPurchaseOrderToSupplyTransactionComponent(
				tagDataSource,
				taggingDataSource,
				memberAccessDataSource,
				purchaseOrderToSupplyDataSource,
			)
			return proposeUpdatepurchaseOrderToSupplyComponent
		},
	)

	container.Transient(
		func(
			purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
			processPOToSupplyComponent purchaseordertosupplydomainrepositoryinterfaces.ProcessPurchaseOrderToSupplyTransactionComponent,
			createNotifComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) purchaseordertosupplydomainrepositoryinterfaces.ProcessPurchaseOrderToSupplyRepository {
			proposeUpdatepurchaseOrderToSupplyRepo, _ := purchaseordertosupplydomainrepositories.NewProcessPurchaseOrderToSupplyRepository(
				purchaseOrderToSupplyDataSource,
				processPOToSupplyComponent,
				createNotifComponent,
				mongoDBTransaction,
			)
			return proposeUpdatepurchaseOrderToSupplyRepo
		},
	)
}
