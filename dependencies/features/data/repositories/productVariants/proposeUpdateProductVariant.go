package productvariantdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	productvariantdomainrepositories "github.com/horeekaa/backend/features/productVariants/data/repositories"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
)

type ProposeUpdateProductVariantDependency struct{}

func (_ *ProposeUpdateProductVariantDependency) Bind() {
	container.Singleton(
		func(
			productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
			updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) productvariantdomainrepositoryinterfaces.ProposeUpdateProductVariantTransactionComponent {
			updateProductVariantComponent, _ := productvariantdomainrepositories.NewProposeUpdateProductVariantTransactionComponent(
				productVariantDataSource,
				loggingDataSource,
				createDescriptivePhotoComponent,
				updateDescriptivePhotoComponent,
				mapProcessorUtility,
			)
			return updateProductVariantComponent
		},
	)
}
