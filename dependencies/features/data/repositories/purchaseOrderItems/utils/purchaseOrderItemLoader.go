package purchaseorderitemdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryutilities "github.com/horeekaa/backend/features/purchaseOrderItems/data/repositories/utils"
	purchaseorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/utils"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
)

type PurchaseOrderItemLoaderDependency struct{}

func (_ *PurchaseOrderItemLoaderDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
			descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
			mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
			productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
			addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
		) purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader {
			purchaseOrderItemLoader, _ := purchaseorderitemdomainrepositoryutilities.NewPurchaseOrderItemLoader(
				accountDataSource,
				personDataSource,
				descriptivePhotoDataSource,
				mouItemDataSource,
				productVariantDataSource,
				productDataSource,
				tagDataSource,
				taggingDataSource,
				addressDataSource,
			)
			return purchaseOrderItemLoader
		},
	)
}
