package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type ManagePersonDeviceTokenDependency struct{}

func (managePrsDeviceTokenDependency *ManagePersonDeviceTokenDependency) bind() {
	container.Singleton(
		func(
			personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
		) accountdomainrepositoryinterfaces.ManagePersonDeviceTokenRepository {
			managePersonDeviceTokenRepository, _ := accountdomainrepositories.NewManagePersonDeviceTokenRepository(
				personDataSource,
			)
			return managePersonDeviceTokenRepository
		},
	)
}
