package addressdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositories "github.com/horeekaa/backend/features/addresses/data/repositories"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	addressdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories/utils"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
)

type ProposeUpdateAddressDependency struct{}

func (_ *ProposeUpdateAddressDependency) Bind() {
	container.Singleton(
		func(
			addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			addressLoader addressdomainrepositoryutilityinterfaces.AddressLoader,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) addressdomainrepositoryinterfaces.ProposeUpdateAddressTransactionComponent {
			updateAddressComponent, _ := addressdomainrepositories.NewProposeUpdateAddressTransactionComponent(
				addressDataSource,
				loggingDataSource,
				addressLoader,
				mapProcessorUtility,
			)
			return updateAddressComponent
		},
	)
}
