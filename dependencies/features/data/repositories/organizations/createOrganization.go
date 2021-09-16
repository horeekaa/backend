package organizationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositories "github.com/horeekaa/backend/features/organizations/data/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
)

type CreateOrganizationDependency struct{}

func (_ *CreateOrganizationDependency) Bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
		) organizationdomainrepositoryinterfaces.CreateOrganizationTransactionComponent {
			createOrganizationComponent, _ := organizationdomainrepositories.NewCreateOrganizationTransactionComponent(
				organizationDataSource,
				loggingDataSource,
			)
			return createOrganizationComponent
		},
	)

	container.Transient(
		func(
			trxComponent organizationdomainrepositoryinterfaces.CreateOrganizationTransactionComponent,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
			bulkCreateTaggingComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) organizationdomainrepositoryinterfaces.CreateOrganizationRepository {
			createOrganizationRepo, _ := organizationdomainrepositories.NewCreateOrganizationRepository(
				trxComponent,
				createDescriptivePhotoComponent,
				bulkCreateTaggingComponent,
				mongoDBTransaction,
			)
			return createOrganizationRepo
		},
	)
}
