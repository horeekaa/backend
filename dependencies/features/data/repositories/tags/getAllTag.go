package tagdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositories "github.com/horeekaa/backend/features/tags/data/repositories"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
)

type GetAllTagDependency struct{}

func (_ *GetAllTagDependency) Bind() {
	container.Singleton(
		func(
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
		) tagdomainrepositoryinterfaces.GetAllTagRepository {
			getAllTagRepo, _ := tagdomainrepositories.NewGetAllTagRepository(
				tagDataSource,
			)
			return getAllTagRepo
		},
	)
}
