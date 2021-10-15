package addressregiongroupdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseaddressregiongroupdatasourceinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/data/dataSources/databases/interfaces/sources"
	addressregiongroupdomainrepositories "github.com/horeekaa/backend/features/addressRegionGroups/data/repositories"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
)

type ApproveUpdateAddressRegionGroupDependency struct{}

func (_ *ApproveUpdateAddressRegionGroupDependency) Bind() {
	container.Singleton(
		func(
			addressRegionGroupDataSource databaseaddressregiongroupdatasourceinterfaces.AddressRegionGroupDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) addressregiongroupdomainrepositoryinterfaces.ApproveUpdateAddressRegionGroupTransactionComponent {
			approveUpdateAddressRegionGroupComponent, _ := addressregiongroupdomainrepositories.NewApproveUpdateAddressRegionGroupTransactionComponent(
				addressRegionGroupDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdateAddressRegionGroupComponent
		},
	)

	container.Transient(
		func(
			trxComponent addressregiongroupdomainrepositoryinterfaces.ApproveUpdateAddressRegionGroupTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) addressregiongroupdomainrepositoryinterfaces.ApproveUpdateAddressRegionGroupRepository {
			approveUpdateAddressRegionGroupRepo, _ := addressregiongroupdomainrepositories.NewApproveUpdateAddressRegionGroupRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return approveUpdateAddressRegionGroupRepo
		},
	)
}
