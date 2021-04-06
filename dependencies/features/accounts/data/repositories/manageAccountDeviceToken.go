package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type ManageAccountDeviceTokenDependency struct{}

func (manageAccDeviceTokenDependency *ManageAccountDeviceTokenDependency) bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
		) accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository {
			manageAccountDeviceTokenRepository, _ := accountdomainrepositories.NewManageAccountDeviceTokenRepository(
				accountDataSource,
			)
			return manageAccountDeviceTokenRepository
		},
	)
}
