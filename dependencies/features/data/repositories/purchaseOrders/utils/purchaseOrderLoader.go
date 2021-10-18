package purchaseorderdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryutilities "github.com/horeekaa/backend/features/purchaseOrders/data/repositories/utils"
	purchaseorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/utils"
)

type PurchaseOrderLoaderDependency struct{}

func (_ *PurchaseOrderLoaderDependency) Bind() {
	container.Singleton(
		func(
			mouDataSource databasemoudatasourceinterfaces.MouDataSource,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
		) purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader {
			purchaseOrderLoader, _ := purchaseorderdomainrepositoryutilities.NewPurchaseOrderLoader(
				mouDataSource,
				organizationDataSource,
				addressDataSource,
			)
			return purchaseOrderLoader
		},
	)
}
