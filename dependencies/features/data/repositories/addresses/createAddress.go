package addressdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositories "github.com/horeekaa/backend/features/addresses/data/repositories"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	addressdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories/utils"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
)

type CreateAddressDependency struct{}

func (_ *CreateAddressDependency) Bind() {
	container.Singleton(
		func(
			addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			addressLoader addressdomainrepositoryutilityinterfaces.AddressLoader,
		) addressdomainrepositoryinterfaces.CreateAddressTransactionComponent {
			createAddressComponent, _ := addressdomainrepositories.NewCreateAddressTransactionComponent(
				addressDataSource,
				loggingDataSource,
				addressLoader,
			)
			return createAddressComponent
		},
	)
}
