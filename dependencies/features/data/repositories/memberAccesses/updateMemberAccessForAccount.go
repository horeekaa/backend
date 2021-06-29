package memberaccessdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositories "github.com/horeekaa/backend/features/memberAccesses/data/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type UpdateMemberAccessForAccountDependency struct{}

func (_ UpdateMemberAccessForAccountDependency) Bind() {
	container.Singleton(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountTransactionComponent {
			updateMemberAccessComponent, _ := memberaccessdomainrepositories.NewUpdateMemberAccessForAccountTransactionComponent(
				memberAccessDataSource,
				memberAccessRefDataSource,
				mapProcessorUtility,
			)
			return updateMemberAccessComponent
		},
	)

	container.Transient(
		func(
			trxComponent memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) memberaccessdomainrepositoryinterfaces.UpdateMemberAccessForAccountRepository {
			updateMemberAccessRepo, _ := memberaccessdomainrepositories.NewUpdateMemberAccessForAccountRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return updateMemberAccessRepo
		},
	)
}
