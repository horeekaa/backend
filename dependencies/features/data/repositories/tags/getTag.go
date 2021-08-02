package tagdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositories "github.com/horeekaa/backend/features/tags/data/repositories"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
)

type GetTagDependency struct{}

func (_ *GetTagDependency) Bind() {
	container.Singleton(
		func(
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
		) tagdomainrepositoryinterfaces.GetTagRepository {
			getTagRepo, _ := tagdomainrepositories.NewGetTagRepository(
				tagDataSource,
			)
			return getTagRepo
		},
	)
}
