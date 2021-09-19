package memberaccessdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositories "github.com/horeekaa/backend/features/memberAccesses/data/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
)

type CreateMemberAccessDependency struct{}

func (createMemberAccessForAccountDependency *CreateMemberAccessDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
		) memberaccessdomainrepositoryinterfaces.CreateMemberAccessTransactionComponent {
			createMemberAccessComponent, _ := memberaccessdomainrepositories.NewCreateMemberAccessTransactionComponent(
				accountDataSource,
				organizationDataSource,
				memberAccessRefDataSource,
				memberAccessDataSource,
				loggingDataSource,
			)
			return createMemberAccessComponent
		},
	)

	container.Transient(
		func(
			createNotifComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
			trxComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository {
			updateMemberAccessRepo, _ := memberaccessdomainrepositories.NewCreateMemberAccessRepository(
				createNotifComponent,
				trxComponent,
				mongoDBTransaction,
			)
			return updateMemberAccessRepo
		},
	)
}
