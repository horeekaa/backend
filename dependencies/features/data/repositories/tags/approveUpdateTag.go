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

type ApproveUpdateTagDependency struct{}

func (_ *ApproveUpdateTagDependency) Bind() {
	container.Singleton(
		func(
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			taggingDataSource databasetaggingdatasourceinterfaces.TaggingDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) tagdomainrepositoryinterfaces.ApproveUpdateTagTransactionComponent {
			approveUpdateTagComponent, _ := tagdomainrepositories.NewApproveUpdateTagTransactionComponent(
				tagDataSource,
				taggingDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdateTagComponent
		},
	)

	container.Transient(
		func(
			trxComponent tagdomainrepositoryinterfaces.ApproveUpdateTagTransactionComponent,
			approveDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent,
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) tagdomainrepositoryinterfaces.ApproveUpdateTagRepository {
			approveUpdateTagRepo, _ := tagdomainrepositories.NewApproveUpdateTagRepository(
				trxComponent,
				approveDescriptivePhotoComponent,
				tagDataSource,
				mongoDBTransaction,
			)
			return approveUpdateTagRepo
		},
	)
}
