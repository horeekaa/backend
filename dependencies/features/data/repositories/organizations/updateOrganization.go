package organizationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositories "github.com/horeekaa/backend/features/organizations/data/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
)

type UpdateOrganizationDependency struct{}

func (_ *UpdateOrganizationDependency) Bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) organizationdomainrepositoryinterfaces.UpdateOrganizationTransactionComponent {
			updateOrganizationComponent, _ := organizationdomainrepositories.NewUpdateOrganizationTransactionComponent(
				organizationDataSource,
				mapProcessorUtility,
			)
			return updateOrganizationComponent
		},
	)

	container.Transient(
		func(
			trxComponent organizationdomainrepositoryinterfaces.UpdateOrganizationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) organizationdomainrepositoryinterfaces.UpdateOrganizationRepository {
			updateOrganizationRepo, _ := organizationdomainrepositories.NewUpdateOrganizationRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return updateOrganizationRepo
		},
	)
}