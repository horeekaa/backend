package addressdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositories "github.com/horeekaa/backend/features/addresses/data/repositories"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
)

type UpdateAddressDependency struct{}

func (_ *UpdateAddressDependency) Bind() {
	container.Singleton(
		func(
			addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
		) addressdomainrepositoryinterfaces.UpdateAddressTransactionComponent {
			updateAddressComponent, _ := addressdomainrepositories.NewUpdateAddressTransactionComponent(
				addressDataSource,
			)
			return updateAddressComponent
		},
	)
}
