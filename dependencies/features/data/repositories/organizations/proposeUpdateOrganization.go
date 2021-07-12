package organizationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositories "github.com/horeekaa/backend/features/organizations/data/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
)

type ProposeUpdateOrganizationDependency struct{}

func (_ *ProposeUpdateOrganizationDependency) Bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
			structComparisonUtility coreutilityinterfaces.StructComparisonUtility,
		) organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationTransactionComponent {
			proposeUpdateOrganizationComponent, _ := organizationdomainrepositories.NewProposeUpdateOrganizationTransactionComponent(
				organizationDataSource,
				loggingDataSource,
				mapProcessorUtility,
				structComparisonUtility,
			)
			return proposeUpdateOrganizationComponent
		},
	)

	container.Transient(
		func(
			trxComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationRepository {
			proposeUpdateOrganizationRepo, _ := organizationdomainrepositories.NewProposeUpdateOrganizationRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return proposeUpdateOrganizationRepo
		},
	)
}
