package supplyorderitemdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryutilities "github.com/horeekaa/backend/features/supplyOrderItems/data/repositories/utils"
	supplyorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/utils"
)

type SupplyOrderItemLoaderDependency struct{}

func (_ *SupplyOrderItemLoaderDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
			addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
		) supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader {
			supplyOrderItemLoader, _ := supplyorderitemdomainrepositoryutilities.NewSupplyOrderItemLoader(
				purchaseOrderToSupplyDataSource,
				addressDataSource,
			)
			return supplyOrderItemLoader
		},
	)
}
