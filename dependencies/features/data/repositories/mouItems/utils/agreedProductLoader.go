package mouitemdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasedescriptivephotodatasourceinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositoryutilities "github.com/horeekaa/backend/features/mouItems/data/repositories/utils"
	mouitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories/utils"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
)

type AgreedProductLoaderDependency struct{}

func (_ *AgreedProductLoaderDependency) Bind() {
	container.Singleton(
		func(
			productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			descriptivePhotoDataSource databasedescriptivephotodatasourceinterfaces.DescriptivePhotoDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader {
			agreedProductLoader, _ := mouitemdomainrepositoryutilities.NewAgreedProductLoader(
				productVariantDataSource,
				productDataSource,
				organizationDataSource,
				descriptivePhotoDataSource,
				mapProcessorUtility,
			)
			return agreedProductLoader
		},
	)
}
