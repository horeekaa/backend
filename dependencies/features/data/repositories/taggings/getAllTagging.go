package taggingdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	taggingdomainrepositories "github.com/horeekaa/backend/features/taggings/data/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
)

type GetAllTaggingDependency struct{}

func (_ *GetAllTaggingDependency) Bind() {
	container.Singleton(
		func(
			taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
		) taggingdomainrepositoryinterfaces.GetAllTaggingRepository {
			getAllTaggingRepo, _ := taggingdomainrepositories.NewGetAllTaggingRepository(
				taggingDataSource,
			)
			return getAllTaggingRepo
		},
	)
}
