package productdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositories "github.com/horeekaa/backend/features/products/data/repositories"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
)

type CreateProductDependency struct{}

func (_ *CreateProductDependency) Bind() {
	container.Singleton(
		func(
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
		) productdomainrepositoryinterfaces.CreateProductTransactionComponent {
			createProductComponent, _ := productdomainrepositories.NewCreateProductTransactionComponent(
				productDataSource,
				loggingDataSource,
				structFieldIteratorUtility,
			)
			return createProductComponent
		},
	)

	container.Transient(
		func(
			createProductComponent productdomainrepositoryinterfaces.CreateProductTransactionComponent,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) productdomainrepositoryinterfaces.CreateProductRepository {
			updateproductRepo, _ := productdomainrepositories.NewCreateProductRepository(
				createProductComponent,
				createDescriptivePhotoComponent,
				mongoDBTransaction,
			)
			return updateproductRepo
		},
	)
}
