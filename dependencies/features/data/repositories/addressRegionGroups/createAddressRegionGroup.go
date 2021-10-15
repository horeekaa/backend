package addressregiongroupdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	addressregiongroupdomainrepositories "github.com/horeekaa/backend/features/addressRegionGroups/data/repositories"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
)

type CreateAddressRegionGroupDependency struct{}

func (_ *CreateAddressRegionGroupDependency) Bind() {
	container.Singleton(
		func(
			addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
		) addressregiongroupdomainrepositoryinterfaces.CreateAddressRegionGroupTransactionComponent {
			createAddressRegionGroupComponent, _ := addressregiongroupdomainrepositories.NewCreateAddressRegionGroupTransactionComponent(
				addressRegionGroupDataSource,
				loggingDataSource,
			)
			return createAddressRegionGroupComponent
		},
	)

	container.Transient(
		func(
			createAddressRegionGroupComponent addressregiongroupdomainrepositoryinterfaces.CreateAddressRegionGroupTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) addressregiongroupdomainrepositoryinterfaces.CreateAddressRegionGroupRepository {
			updateAddressRegionGroupRepo, _ := addressregiongroupdomainrepositories.NewCreateAddressRegionGroupRepository(
				createAddressRegionGroupComponent,
				mongoDBTransaction,
			)
			return updateAddressRegionGroupRepo
		},
	)
}
