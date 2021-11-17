package addressdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositories "github.com/horeekaa/backend/features/addresses/data/repositories"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
)

type ApproveUpdateAddressDependency struct{}

func (_ *ApproveUpdateAddressDependency) Bind() {
	container.Singleton(
		func(
			addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) addressdomainrepositoryinterfaces.ApproveUpdateAddressTransactionComponent {
			approveUpdateAddressComponent, _ := addressdomainrepositories.NewApproveUpdateAddressTransactionComponent(
				addressDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdateAddressComponent
		},
	)
}
