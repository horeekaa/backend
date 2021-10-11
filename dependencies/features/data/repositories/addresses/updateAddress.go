package addressdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositories "github.com/horeekaa/backend/features/addresses/data/repositories"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	addressdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories/utils"
)

type UpdateAddressDependency struct{}

func (_ *UpdateAddressDependency) Bind() {
	container.Singleton(
		func(
			addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
			addressLoader addressdomainrepositoryutilityinterfaces.AddressLoader,
		) addressdomainrepositoryinterfaces.UpdateAddressTransactionComponent {
			updateAddressComponent, _ := addressdomainrepositories.NewUpdateAddressTransactionComponent(
				addressDataSource,
				addressLoader,
			)
			return updateAddressComponent
		},
	)
}
