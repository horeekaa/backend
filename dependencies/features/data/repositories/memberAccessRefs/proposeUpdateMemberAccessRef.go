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

type ProposeUpdateMemberAccessRefDependency struct{}

func (_ *ProposeUpdateMemberAccessRefDependency) Bind() {
	container.Singleton(
		func(
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefTransactionComponent {
			proposeUpdateMemberAccessRefComponent, _ := memberaccessrefdomainrepositories.NewProposeUpdateMemberAccessRefTransactionComponent(
				memberAccessRefDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return proposeUpdateMemberAccessRefComponent
		},
	)

	container.Transient(
		func(
			trxComponent memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) memberaccessrefdomainrepositoryinterfaces.ProposeUpdateMemberAccessRefRepository {
			proposeUpdateMemberAccessRefRepo, _ := memberaccessrefdomainrepositories.NewProposeUpdateMemberAccessRefRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return proposeUpdateMemberAccessRefRepo
		},
	)
}
