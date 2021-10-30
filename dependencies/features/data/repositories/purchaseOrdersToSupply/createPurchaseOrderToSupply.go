package purchaseordertosupplydomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	purchaseordertosupplydomainrepositories "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/repositories"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
)

type CreatePurchaseOrderToSupplyDependency struct{}

func (_ *CreatePurchaseOrderToSupplyDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
			purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
			purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
		) purchaseordertosupplydomainrepositoryinterfaces.CreatePurchaseOrderToSupplyTransactionComponent {
			createPurchaseOrderToSupplyComponent, _ := purchaseordertosupplydomainrepositories.NewCreatePurchaseOrderToSupplyTransactionComponent(
				purchaseOrderDataSource,
				purchaseOrderItemDataSource,
				purchaseOrderToSupplyDataSource,
			)
			return createPurchaseOrderToSupplyComponent
		},
	)

	container.Transient(
		func(
			purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
			createPOToSupplyComponent purchaseordertosupplydomainrepositoryinterfaces.CreatePurchaseOrderToSupplyTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) purchaseordertosupplydomainrepositoryinterfaces.CreatePurchaseOrderToSupplyRepository {
			createPurchaseOrderToSupplyRepo, _ := purchaseordertosupplydomainrepositories.NewCreatePurchaseOrderToSupplyRepository(
				purchaseOrderDataSource,
				createPOToSupplyComponent,
				mongoDBTransaction,
			)
			return createPurchaseOrderToSupplyRepo
		},
	)
}
