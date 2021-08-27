package memberaccessdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositories "github.com/horeekaa/backend/features/memberAccesses/data/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
)

type ApproveUpdateMemberAccessDependency struct{}

func (_ *ApproveUpdateMemberAccessDependency) Bind() {
	container.Singleton(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessTransactionComponent {
			approveUpdateMemberAccessComponent, _ := memberaccessdomainrepositories.NewApproveUpdateMemberAccessTransactionComponent(
				memberAccessDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdateMemberAccessComponent
		},
	)

	container.Transient(
		func(
			createNotifComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
			trxComponent memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) memberaccessdomainrepositoryinterfaces.ApproveUpdateMemberAccessRepository {
			approveUpdateMemberAccessRepo, _ := memberaccessdomainrepositories.NewApproveUpdateMemberAccessRepository(
				createNotifComponent,
				trxComponent,
				mongoDBTransaction,
			)
			return approveUpdateMemberAccessRepo
		},
	)
}
