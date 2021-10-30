package purchaseordertosupplydomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	purchaseordertosupplydomainrepositories "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/repositories"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
)

type GetAllPurchaseOrderToSupplyDependency struct{}

func (_ *GetAllPurchaseOrderToSupplyDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) purchaseordertosupplydomainrepositoryinterfaces.GetAllPurchaseOrderToSupplyRepository {
			getAllPurchaseOrderToSupplyRepo, _ := purchaseordertosupplydomainrepositories.NewGetAllPurchaseOrderToSupplyRepository(
				purchaseOrderToSupplyDataSource,
				mongoQueryBuilder,
			)
			return getAllPurchaseOrderToSupplyRepo
		},
	)
}
