package tagdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasetagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/interfaces/sources"
	tagdomainrepositories "github.com/horeekaa/backend/features/tags/data/repositories"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
)

type ApproveUpdateTagDependency struct{}

func (_ *ApproveUpdateTagDependency) Bind() {
	container.Singleton(
		func(
			tagDataSource databasetagdatasourceinterfaces.TagDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) tagdomainrepositoryinterfaces.ApproveUpdateTagTransactionComponent {
			approveUpdateTagComponent, _ := tagdomainrepositories.NewApproveUpdateTagTransactionComponent(
				tagDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdateTagComponent
		},
	)

	container.Transient(
		func(
			trxComponent tagdomainrepositoryinterfaces.ApproveUpdateTagTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) tagdomainrepositoryinterfaces.ApproveUpdateTagRepository {
			approveUpdateTagRepo, _ := tagdomainrepositories.NewApproveUpdateTagRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return approveUpdateTagRepo
		},
	)
}
