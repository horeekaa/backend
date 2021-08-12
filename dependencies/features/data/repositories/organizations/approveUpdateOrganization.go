package organizationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositories "github.com/horeekaa/backend/features/organizations/data/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
)

type ApproveUpdateOrganizationDependency struct{}

func (_ *ApproveUpdateOrganizationDependency) Bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationTransactionComponent {
			approveUpdateOrganizationComponent, _ := organizationdomainrepositories.NewApproveUpdateOrganizationTransactionComponent(
				organizationDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdateOrganizationComponent
		},
	)

	container.Transient(
		func(
			trxComponent organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationTransactionComponent,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			bulkApproveUpdateTaggingComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationRepository {
			approveUpdateOrganizationRepo, _ := organizationdomainrepositories.NewApproveUpdateOrganizationRepository(
				trxComponent,
				organizationDataSource,
				bulkApproveUpdateTaggingComponent,
				mongoDBTransaction,
			)
			return approveUpdateOrganizationRepo
		},
	)
}
