package taggingdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	taggingdomainrepositories "github.com/horeekaa/backend/features/taggings/data/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
)

type GetTaggingDependency struct{}

func (_ *GetTaggingDependency) Bind() {
	container.Singleton(
		func(
			taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
		) taggingdomainrepositoryinterfaces.GetTaggingRepository {
			getTaggingRepo, _ := taggingdomainrepositories.NewGetTaggingRepository(
				taggingDataSource,
			)
			return getTaggingRepo
		},
	)
}
