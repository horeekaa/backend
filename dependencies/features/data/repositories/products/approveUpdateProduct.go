package productdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositories "github.com/horeekaa/backend/features/products/data/repositories"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
)

type ApproveUpdateProductDependency struct{}

func (_ *ApproveUpdateProductDependency) Bind() {
	container.Singleton(
		func(
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) productdomainrepositoryinterfaces.ApproveUpdateProductTransactionComponent {
			approveUpdateProductComponent, _ := productdomainrepositories.NewApproveUpdateProductTransactionComponent(
				productDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdateProductComponent
		},
	)

	container.Transient(
		func(
			trxComponent productdomainrepositoryinterfaces.ApproveUpdateProductTransactionComponent,
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
			approveUpdateProductVariantComponent productvariantdomainrepositoryinterfaces.ApproveUpdateProductVariantTransactionComponent,
			bulkApproveUpdateTaggingComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) productdomainrepositoryinterfaces.ApproveUpdateProductRepository {
			approveUpdateProductRepo, _ := productdomainrepositories.NewApproveUpdateProductRepository(
				trxComponent,
				productDataSource,
				approveUpdateProductVariantComponent,
				bulkApproveUpdateTaggingComponent,
				mongoDBTransaction,
			)
			return approveUpdateProductRepo
		},
	)
}
