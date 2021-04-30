package organizationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositories "github.com/horeekaa/backend/features/organizations/data/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
)

type UpdateOrganizationDependency struct{}

func (_ *UpdateOrganizationDependency) bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
		) organizationdomainrepositoryinterfaces.UpdateOrganizationTransactionComponent {
			updateOrganizationComponent, _ := organizationdomainrepositories.NewUpdateOrganizationTransactionComponent(
				organizationDataSource,
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
