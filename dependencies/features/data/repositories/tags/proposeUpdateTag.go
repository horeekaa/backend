package tagdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasetaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/interfaces/sources"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositories "github.com/horeekaa/backend/features/tags/data/repositories"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
)

type ProposeUpdateTagDependency struct{}

func (_ *ProposeUpdateTagDependency) Bind() {
	container.Singleton(
		func(
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) tagdomainrepositoryinterfaces.ProposeUpdateTagTransactionComponent {
			proposeUpdateTagComponent, _ := tagdomainrepositories.NewProposeUpdateTagTransactionComponent(
				tagDataSource,
				taggingDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return proposeUpdateTagComponent
		},
	)

	container.Transient(
		func(
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
			updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent,
			proposeUpdateTagComponent tagdomainrepositoryinterfaces.ProposeUpdateTagTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) tagdomainrepositoryinterfaces.ProposeUpdateTagRepository {
			proposeUpdateTagRepo, _ := tagdomainrepositories.NewProposeUpdateTagRepository(
				tagDataSource,
				createDescriptivePhotoComponent,
				updateDescriptivePhotoComponent,
				proposeUpdateTagComponent,
				mongoDBTransaction,
			)
			return proposeUpdateTagRepo
		},
	)
}
