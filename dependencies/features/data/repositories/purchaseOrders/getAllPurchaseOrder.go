package purchaseorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositories "github.com/horeekaa/backend/features/purchaseOrders/data/repositories"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
)

type GetAllPurchaseOrderDependency struct{}

func (_ *GetAllPurchaseOrderDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) purchaseorderdomainrepositoryinterfaces.GetAllPurchaseOrderRepository {
			getAllPurchaseOrderRepo, _ := purchaseorderdomainrepositories.NewGetAllPurchaseOrderRepository(
				purchaseOrderDataSource,
				mongoQueryBuilder,
			)
			return getAllPurchaseOrderRepo
		},
	)
}
