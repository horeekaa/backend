package purchaseorderitemdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryutilities "github.com/horeekaa/backend/features/purchaseOrderItems/data/repositories/utils"
	purchaseorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/utils"
)

type PurchaseOrderItemLoaderDependency struct{}

func (_ *PurchaseOrderItemLoaderDependency) Bind() {
	container.Singleton(
		func(
			descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
			mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
			productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
		) purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader {
			purchaseOrderItemLoader, _ := purchaseorderitemdomainrepositoryutilities.NewPurchaseOrderItemLoader(
				descriptivePhotoDataSource,
				mouItemDataSource,
				productVariantDataSource,
				productDataSource,
			)
			return purchaseOrderItemLoader
		},
	)
}
