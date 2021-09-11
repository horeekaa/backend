package mouitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasemouItemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositories "github.com/horeekaa/backend/features/mouItems/data/repositories"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
)

type GetMouItemDependency struct{}

func (_ *GetMouItemDependency) Bind() {
	container.Singleton(
		func(
			mouItemDataSource databasemouItemdatasourceinterfaces.MouItemDataSource,
		) mouitemdomainrepositoryinterfaces.GetMouItemRepository {
			getMouItemRepo, _ := mouitemdomainrepositories.NewGetMouItemRepository(
				mouItemDataSource,
			)
			return getMouItemRepo
		},
	)
}
