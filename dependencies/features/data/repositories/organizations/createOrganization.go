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

type CreateOrganizationDependency struct{}

func (_ *CreateOrganizationDependency) Bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
		) organizationdomainrepositoryinterfaces.CreateOrganizationTransactionComponent {
			createOrganizationComponent, _ := organizationdomainrepositories.NewCreateOrganizationTransactionComponent(
				organizationDataSource,
				loggingDataSource,
				structFieldIteratorUtility,
			)
			return createOrganizationComponent
		},
	)

	container.Transient(
		func(
			trxComponent organizationdomainrepositoryinterfaces.CreateOrganizationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) organizationdomainrepositoryinterfaces.CreateOrganizationRepository {
			createOrganizationRepo, _ := organizationdomainrepositories.NewCreateOrganizationRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return createOrganizationRepo
		},
	)
}
