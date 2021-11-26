package supplyorderdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryutilities "github.com/horeekaa/backend/features/supplyOrders/data/repositories/utils"
	supplyorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/utils"
)

type SupplyOrderLoaderDependency struct{}

func (_ *SupplyOrderLoaderDependency) Bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
		) supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader {
			supplyOrderLoader, _ := supplyorderdomainrepositoryutilities.NewSupplyOrderLoader(
				organizationDataSource,
			)
			return supplyOrderLoader
		},
	)
}
