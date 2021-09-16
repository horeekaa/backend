package productdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositories "github.com/horeekaa/backend/features/products/data/repositories"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
)

type ProposeUpdateProductDependency struct{}

func (_ *ProposeUpdateProductDependency) Bind() {
	container.Singleton(
		func(
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) productdomainrepositoryinterfaces.ProposeUpdateProductTransactionComponent {
			proposeUpdateProductComponent, _ := productdomainrepositories.NewProposeUpdateProductTransactionComponent(
				productDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return proposeUpdateProductComponent
		},
	)

	container.Transient(
		func(
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
			updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent,
			createProductVariantComponent productvariantdomainrepositoryinterfaces.CreateProductVariantTransactionComponent,
			updateProductVariantComponent productvariantdomainrepositoryinterfaces.UpdateProductVariantTransactionComponent,
			createTaggingComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent,
			updateTaggingComponent taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingTransactionComponent,
			proposeUpdateProductComponent productdomainrepositoryinterfaces.ProposeUpdateProductTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) productdomainrepositoryinterfaces.ProposeUpdateProductRepository {
			proposeUpdateProductRepo, _ := productdomainrepositories.NewProposeUpdateProductRepository(
				productDataSource,
				createDescriptivePhotoComponent,
				updateDescriptivePhotoComponent,
				createProductVariantComponent,
				updateProductVariantComponent,
				createTaggingComponent,
				updateTaggingComponent,
				proposeUpdateProductComponent,
				mongoDBTransaction,
			)
			return proposeUpdateProductRepo
		},
	)
}
