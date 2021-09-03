package moudomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositories "github.com/horeekaa/backend/features/mous/data/repositories"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
)

type ApproveUpdateMouDependency struct{}

func (_ *ApproveUpdateMouDependency) Bind() {
	container.Singleton(
		func(
			mouDataSource databasemoudatasourceinterfaces.MouDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) moudomainrepositoryinterfaces.ApproveUpdateMouTransactionComponent {
			approveUpdateMouComponent, _ := moudomainrepositories.NewApproveUpdateMouTransactionComponent(
				mouDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdateMouComponent
		},
	)

	container.Transient(
		func(
			trxComponent moudomainrepositoryinterfaces.ApproveUpdateMouTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) moudomainrepositoryinterfaces.ApproveUpdateMouRepository {
			approveUpdateMouRepo, _ := moudomainrepositories.NewApproveUpdateMouRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return approveUpdateMouRepo
		},
	)
}
