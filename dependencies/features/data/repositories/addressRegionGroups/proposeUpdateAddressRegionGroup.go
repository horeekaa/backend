package addressregiongroupdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseaddressRegionGroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	addressregiongroupdomainrepositories "github.com/horeekaa/backend/features/addressRegionGroups/data/repositories"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
)

type ProposeUpdateAddressRegionGroupDependency struct{}

func (_ *ProposeUpdateAddressRegionGroupDependency) Bind() {
	container.Singleton(
		func(
			addressRegionGroupDataSource databaseaddressRegionGroupdatasourceinterfaces.AddressRegionGroupDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) addressregiongroupdomainrepositoryinterfaces.ProposeUpdateAddressRegionGroupTransactionComponent {
			proposeUpdateAddressRegionGroupComponent, _ := addressregiongroupdomainrepositories.NewProposeUpdateAddressRegionGroupTransactionComponent(
				addressRegionGroupDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return proposeUpdateAddressRegionGroupComponent
		},
	)

	container.Transient(
		func(
			proposeUpdateAddressRegionGroupComponent addressregiongroupdomainrepositoryinterfaces.ProposeUpdateAddressRegionGroupTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) addressregiongroupdomainrepositoryinterfaces.ProposeUpdateAddressRegionGroupRepository {
			proposeUpdateAddressRegionGroupRepo, _ := addressregiongroupdomainrepositories.NewProposeUpdateAddressRegionGroupRepository(
				proposeUpdateAddressRegionGroupComponent,
				mongoDBTransaction,
			)
			return proposeUpdateAddressRegionGroupRepo
		},
	)
}
