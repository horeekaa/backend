package addressdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	addressdomainrepositories "github.com/horeekaa/backend/features/addresses/data/repositories"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
)

type GetAddressDependency struct{}

func (_ *GetAddressDependency) Bind() {
	container.Singleton(
		func(
			addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
		) addressdomainrepositoryinterfaces.GetAddressRepository {
			getAddressRepo, _ := addressdomainrepositories.NewGetAddressRepository(
				addressDataSource,
			)
			return getAddressRepo
		},
	)
}
