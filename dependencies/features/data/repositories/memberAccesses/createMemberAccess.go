package memberaccessdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositories "github.com/horeekaa/backend/features/memberAccesses/data/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
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
			structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
		) memberaccessdomainrepositoryinterfaces.CreateMemberAccessTransactionComponent {
			createMemberAccessComponent, _ := memberaccessdomainrepositories.NewCreateMemberAccessTransactionComponent(
				accountDataSource,
				organizationDataSource,
				memberAccessRefDataSource,
				memberAccessDataSource,
				loggingDataSource,
				structFieldIteratorUtility,
			)
			return createMemberAccessComponent
		},
	)

	container.Transient(
		func(
			trxComponent memberaccessdomainrepositoryinterfaces.CreateMemberAccessTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository {
			updateMemberAccessRepo, _ := memberaccessdomainrepositories.NewCreateMemberAccessRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return updateMemberAccessRepo
		},
	)
}
