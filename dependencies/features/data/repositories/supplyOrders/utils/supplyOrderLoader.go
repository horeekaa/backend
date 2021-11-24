package supplyorderdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryutilities "github.com/horeekaa/backend/features/supplyOrders/data/repositories/utils"
	supplyorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/utils"
)

type SupplyOrderLoaderDependency struct{}

func (_ *SupplyOrderLoaderDependency) Bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
		) supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader {
			supplyOrderLoader, _ := supplyorderdomainrepositoryutilities.NewSupplyOrderLoader(
				organizationDataSource,
				addressDataSource,
			)
			return supplyOrderLoader
		},
	)
}
