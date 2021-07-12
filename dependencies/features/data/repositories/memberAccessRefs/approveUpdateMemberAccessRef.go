package memberaccessrefdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositories "github.com/horeekaa/backend/features/memberAccessRefs/data/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
)

type ApproveUpdateMemberAccessRefDependency struct{}

func (_ *ApproveUpdateMemberAccessRefDependency) Bind() {
	container.Singleton(
		func(
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) memberaccessrefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefTransactionComponent {
			approveUpdateMemberAccessRefComponent, _ := memberaccessrefdomainrepositories.NewApproveUpdateMemberAccessRefTransactionComponent(
				memberAccessRefDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdateMemberAccessRefComponent
		},
	)

	container.Transient(
		func(
			trxComponent memberaccessrefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) memberaccessrefdomainrepositoryinterfaces.ApproveUpdateMemberAccessRefRepository {
			approveUpdateMemberAccessRefRepo, _ := memberaccessrefdomainrepositories.NewApproveUpdateMemberAccessRefRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return approveUpdateMemberAccessRefRepo
		},
	)
}
