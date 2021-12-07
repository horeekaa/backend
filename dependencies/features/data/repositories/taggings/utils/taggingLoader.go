package taggingitemdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	taggingdomainrepositoryutilities "github.com/horeekaa/backend/features/taggings/data/repositories/utils"
	taggingdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories/utils"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
)

type TaggingLoaderDependency struct{}

func (_ *TaggingLoaderDependency) Bind() {
	container.Singleton(
		func(
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
		) taggingdomainrepositoryutilityinterfaces.TaggingLoader {
			taggingItemLoader, _ := taggingdomainrepositoryutilities.NewTaggingLoader(
				tagDataSource,
			)
			return taggingItemLoader
		},
	)
}
