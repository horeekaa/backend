package mouitemdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	mouitemdomainrepositoryutilities "github.com/horeekaa/backend/features/mouItems/data/repositories/utils"
	mouitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories/utils"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
)

type AgreedProductLoaderDependency struct{}

func (_ *AgreedProductLoaderDependency) Bind() {
	container.Singleton(
		func(
			productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader {
			agreedProductLoader, _ := mouitemdomainrepositoryutilities.NewAgreedProductLoader(
				productVariantDataSource,
				productDataSource,
				mapProcessorUtility,
			)
			return agreedProductLoader
		},
	)
}
