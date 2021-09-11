package moudomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositories "github.com/horeekaa/backend/features/mous/data/repositories"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
)

type GetAllMouDependency struct{}

func (_ *GetAllMouDependency) Bind() {
	container.Singleton(
		func(
			mouDataSource databasemoudatasourceinterfaces.MouDataSource,
		) moudomainrepositoryinterfaces.GetAllMouRepository {
			getAllMouRepo, _ := moudomainrepositories.NewGetAllMouRepository(
				mouDataSource,
			)
			return getAllMouRepo
		},
	)
}
