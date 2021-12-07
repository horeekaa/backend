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
	taggingdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories/utils"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
)

type BulkProposeUpdateTaggingDependency struct{}

func (_ *BulkProposeUpdateTaggingDependency) Bind() {
	container.Singleton(
		func(
			taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
			taggingLoaderUtility taggingdomainrepositoryutilityinterfaces.TaggingLoader,
		) taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingTransactionComponent {
			bulkProposeUpdateTaggingComponent, _ := taggingdomainrepositories.NewBulkProposeUpdateTaggingTransactionComponent(
				taggingDataSource,
				tagDataSource,
				organizationDataSource,
				productDataSource,
				loggingDataSource,
				mapProcessorUtility,
				taggingLoaderUtility,
			)
			return bulkProposeUpdateTaggingComponent
		},
	)

	container.Transient(
		func(
			bulkProposeUpdateTaggingComponent taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingRepository {
			bulkProposeUpdateTaggingRepo, _ := taggingdomainrepositories.NewBulkProposeUpdateTaggingRepository(
				bulkProposeUpdateTaggingComponent,
				mongoDBTransaction,
			)
			return bulkProposeUpdateTaggingRepo
		},
	)
}
