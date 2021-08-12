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

type BulkCreateTaggingDependency struct{}

func (_ *BulkCreateTaggingDependency) Bind() {
	container.Singleton(
		func(
			taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
			structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
		) taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent {
			bulkCreateTaggingComponent, _ := taggingdomainrepositories.NewBulkCreateTaggingTransactionComponent(
				taggingDataSource,
				loggingDataSource,
				tagDataSource,
				organizationDataSource,
				productDataSource,
				structFieldIteratorUtility,
			)
			return bulkCreateTaggingComponent
		},
	)

	container.Transient(
		func(
			bulkCreateTaggingComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) taggingdomainrepositoryinterfaces.BulkCreateTaggingRepository {
			bulkCreateTaggingRepo, _ := taggingdomainrepositories.NewBulkCreateTaggingRepository(
				bulkCreateTaggingComponent,
				mongoDBTransaction,
			)
			return bulkCreateTaggingRepo
		},
	)
}
