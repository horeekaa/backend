package tagdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositories "github.com/horeekaa/backend/features/tags/data/repositories"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
)

type CreateTagDependency struct{}

func (_ *CreateTagDependency) Bind() {
	container.Singleton(
		func(
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
		) tagdomainrepositoryinterfaces.CreateTagTransactionComponent {
			createTagComponent, _ := tagdomainrepositories.NewCreateTagTransactionComponent(
				tagDataSource,
				loggingDataSource,
			)
			return createTagComponent
		},
	)

	container.Transient(
		func(
			createTagComponent tagdomainrepositoryinterfaces.CreateTagTransactionComponent,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) tagdomainrepositoryinterfaces.CreateTagRepository {
			updateTagRepo, _ := tagdomainrepositories.NewCreateTagRepository(
				createTagComponent,
				createDescriptivePhotoComponent,
				mongoDBTransaction,
			)
			return updateTagRepo
		},
	)
}
