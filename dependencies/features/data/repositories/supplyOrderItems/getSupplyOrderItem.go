package supplyorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositories "github.com/horeekaa/backend/features/supplyOrderItems/data/repositories"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
)

type GetSupplyOrderItemDependency struct{}

func (_ *GetSupplyOrderItemDependency) Bind() {
	container.Singleton(
		func(
			supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
		) supplyorderitemdomainrepositoryinterfaces.GetSupplyOrderItemRepository {
			getSupplyOrderItemRepo, _ := supplyorderitemdomainrepositories.NewGetSupplyOrderItemRepository(
				supplyOrderItemDataSource,
			)
			return getSupplyOrderItemRepo
		},
	)
}
