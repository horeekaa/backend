package productvariantdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	productvariantdomainrepositories "github.com/horeekaa/backend/features/productVariants/data/repositories"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
)

type CreateProductVariantDependency struct{}

func (_ *CreateProductVariantDependency) Bind() {
	container.Singleton(
		func(
			productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
		) productvariantdomainrepositoryinterfaces.CreateProductVariantTransactionComponent {
			createProductVariantComponent, _ := productvariantdomainrepositories.NewCreateProductVariantTransactionComponent(
				productVariantDataSource,
				loggingDataSource,
				createDescriptivePhotoComponent,
			)
			return createProductVariantComponent
		},
	)
}
