package purchaseorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositories "github.com/horeekaa/backend/features/purchaseOrderItems/data/repositories"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
)

type GetAllPurchaseOrderItemDependency struct{}

func (_ *GetAllPurchaseOrderItemDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) purchaseorderitemdomainrepositoryinterfaces.GetAllPurchaseOrderItemRepository {
			getAllPurchaseOrderItemRepo, _ := purchaseorderitemdomainrepositories.NewGetAllPurchaseOrderItemRepository(
				purchaseOrderItemDataSource,
				mongoQueryBuilder,
			)
			return getAllPurchaseOrderItemRepo
		},
	)
}
