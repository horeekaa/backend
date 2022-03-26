package moudomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositories "github.com/horeekaa/backend/features/mous/data/repositories"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moudomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories/utils"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
)

type CreateMouDependency struct{}

func (_ *CreateMouDependency) Bind() {
	container.Singleton(
		func(
			mouDataSource databasemoudatasourceinterfaces.MouDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			partyLoader moudomainrepositoryutilityinterfaces.PartyLoader,
		) moudomainrepositoryinterfaces.CreateMouTransactionComponent {
			createmouComponent, _ := moudomainrepositories.NewCreateMouTransactionComponent(
				mouDataSource,
				loggingDataSource,
				partyLoader,
			)
			return createmouComponent
		},
	)

	container.Transient(
		func(
			trxComponent moudomainrepositoryinterfaces.CreateMouTransactionComponent,
			createMouItemComponent mouitemdomainrepositoryinterfaces.CreateMouItemTransactionComponent,
			createNotifComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) moudomainrepositoryinterfaces.CreateMouRepository {
			createMouRepo, _ := moudomainrepositories.NewCreateMouRepository(
				trxComponent,
				createMouItemComponent,
				createNotifComponent,
				mongoDBTransaction,
			)
			return createMouRepo
		},
	)
}
