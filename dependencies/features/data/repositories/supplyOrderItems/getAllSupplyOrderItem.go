package supplyorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositories "github.com/horeekaa/backend/features/supplyOrderItems/data/repositories"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
)

type GetAllSupplyOrderItemDependency struct{}

func (_ *GetAllSupplyOrderItemDependency) Bind() {
	container.Singleton(
		func(
			supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) supplyorderitemdomainrepositoryinterfaces.GetAllSupplyOrderItemRepository {
			getAllSupplyOrderItemRepo, _ := supplyorderitemdomainrepositories.NewGetAllSupplyOrderItemRepository(
				supplyOrderItemDataSource,
				mongoQueryBuilder,
			)
			return getAllSupplyOrderItemRepo
		},
	)
}
