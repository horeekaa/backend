package purchaseorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databasepurchaseorderItemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositories "github.com/horeekaa/backend/features/purchaseOrderItems/data/repositories"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
)

type ProposeUpdatePurchaseOrderItemDeliveryDependency struct{}

func (_ *ProposeUpdatePurchaseOrderItemDeliveryDependency) Bind() {
	container.Transient(
		func(
			purchaseOrderItemDataSource databasepurchaseorderItemdatasourceinterfaces.PurchaseOrderItemDataSource,
			proposeUpdatePurchaseOrderItemTransactionComponent purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemDeliveryRepository {
			proposeUpdatePurchaseOrderItemDeliveryRepo, _ := purchaseorderitemdomainrepositories.NewProposeUpdatePurchaseOrderItemDeliveryRepository(
				purchaseOrderItemDataSource,
				proposeUpdatePurchaseOrderItemTransactionComponent,
				mongoDBTransaction,
			)
			return proposeUpdatePurchaseOrderItemDeliveryRepo
		},
	)
}
