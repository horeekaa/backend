package supplyorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositories "github.com/horeekaa/backend/features/supplyOrders/data/repositories"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
)

type GetAllSupplyOrderDependency struct{}

func (_ *GetAllSupplyOrderDependency) Bind() {
	container.Singleton(
		func(
			supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) supplyorderdomainrepositoryinterfaces.GetAllSupplyOrderRepository {
			getAllSupplyOrderRepo, _ := supplyorderdomainrepositories.NewGetAllSupplyOrderRepository(
				supplyOrderDataSource,
				mongoQueryBuilder,
			)
			return getAllSupplyOrderRepo
		},
	)
}
