package taggingdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	taggingdomainrepositories "github.com/horeekaa/backend/features/taggings/data/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
)

type BulkApproveUpdateTaggingDependency struct{}

func (_ *BulkApproveUpdateTaggingDependency) Bind() {
	container.Singleton(
		func(
			taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent {
			bulkApproveUpdateTaggingComponent, _ := taggingdomainrepositories.NewBulkApproveUpdateTaggingTransactionComponent(
				taggingDataSource,
				tagDataSource,
				organizationDataSource,
				productDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return bulkApproveUpdateTaggingComponent
		},
	)

	container.Transient(
		func(
			bulkApproveUpdateTaggingComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingRepository {
			bulkApproveUpdateTaggingRepo, _ := taggingdomainrepositories.NewBulkApproveUpdateTaggingRepository(
				bulkApproveUpdateTaggingComponent,
				mongoDBTransaction,
			)
			return bulkApproveUpdateTaggingRepo
		},
	)
}
