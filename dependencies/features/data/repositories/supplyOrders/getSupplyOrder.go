package supplyorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasesupplyOrderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositories "github.com/horeekaa/backend/features/supplyOrders/data/repositories"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
)

type GetSupplyOrderDependency struct{}

func (_ *GetSupplyOrderDependency) Bind() {
	container.Singleton(
		func(
			supplyOrderDataSource databasesupplyOrderdatasourceinterfaces.SupplyOrderDataSource,
		) supplyorderdomainrepositoryinterfaces.GetSupplyOrderRepository {
			getsupplyOrderRepo, _ := supplyorderdomainrepositories.NewGetSupplyOrderRepository(
				supplyOrderDataSource,
			)
			return getsupplyOrderRepo
		},
	)
}
