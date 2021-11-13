package organizationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositories "github.com/horeekaa/backend/features/organizations/data/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
)

type ProposeUpdateOrganizationDependency struct{}

func (_ *ProposeUpdateOrganizationDependency) Bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationTransactionComponent {
			proposeUpdateOrganizationComponent, _ := organizationdomainrepositories.NewProposeUpdateOrganizationTransactionComponent(
				organizationDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return proposeUpdateOrganizationComponent
		},
	)

	container.Transient(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
			proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent,
			createAddressComponent addressdomainrepositoryinterfaces.CreateAddressTransactionComponent,
			updateAddressComponent addressdomainrepositoryinterfaces.UpdateAddressTransactionComponent,
			bulkCreateTaggingComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent,
			bulkUpdateTaggingComponent taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingTransactionComponent,
			trxComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationRepository {
			proposeUpdateOrganizationRepo, _ := organizationdomainrepositories.NewProposeUpdateOrganizationRepository(
				organizationDataSource,
				createDescriptivePhotoComponent,
				proposeUpdateDescriptivePhotoComponent,
				createAddressComponent,
				updateAddressComponent,
				bulkCreateTaggingComponent,
				bulkUpdateTaggingComponent,
				trxComponent,
				mongoDBTransaction,
			)
			return proposeUpdateOrganizationRepo
		},
	)
}
