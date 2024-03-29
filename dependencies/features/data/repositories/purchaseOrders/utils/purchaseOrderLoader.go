package purchaseorderdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
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
		) purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader {
			purchaseOrderLoader, _ := purchaseorderdomainrepositoryutilities.NewPurchaseOrderLoader(
				mouDataSource,
				organizationDataSource,
			)
			return purchaseOrderLoader
		},
	)
}
