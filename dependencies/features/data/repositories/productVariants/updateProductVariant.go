package productvariantdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	productvariantdomainrepositories "github.com/horeekaa/backend/features/productVariants/data/repositories"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
)

type UpdateProductVariantDependency struct{}

func (_ *UpdateProductVariantDependency) Bind() {
	container.Singleton(
		func(
			productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
			updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent,
		) productvariantdomainrepositoryinterfaces.UpdateProductVariantTransactionComponent {
			updateProductVariantComponent, _ := productvariantdomainrepositories.NewUpdateProductVariantTransactionComponent(
				productVariantDataSource,
				createDescriptivePhotoComponent,
				updateDescriptivePhotoComponent,
			)
			return updateProductVariantComponent
		},
	)
}
