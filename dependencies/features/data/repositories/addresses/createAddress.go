package addressdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositories "github.com/horeekaa/backend/features/addresses/data/repositories"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	addressdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories/utils"
)

type CreateAddressDependency struct{}

func (_ *CreateAddressDependency) Bind() {
	container.Singleton(
		func(
			addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
			addressLoader addressdomainrepositoryutilityinterfaces.AddressLoader,
		) addressdomainrepositoryinterfaces.CreateAddressTransactionComponent {
			createAddressComponent, _ := addressdomainrepositories.NewCreateAddressTransactionComponent(
				addressDataSource,
				addressLoader,
			)
			return createAddressComponent
		},
	)
}
