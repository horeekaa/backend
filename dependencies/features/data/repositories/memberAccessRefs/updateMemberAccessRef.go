package memberaccessrefdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositories "github.com/horeekaa/backend/features/memberAccessRefs/data/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
)

type UpdateMemberAccessRefDependency struct{}

func (_ UpdateMemberAccessRefDependency) Bind() {
	container.Singleton(
		func(
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefTransactionComponent {
			updateMemberAccessRefComponent, _ := memberaccessrefdomainrepositories.NewUpdateMemberAccessRefTransactionComponent(
				memberAccessRefDataSource,
				mapProcessorUtility,
			)
			return updateMemberAccessRefComponent
		},
	)

	container.Transient(
		func(
			trxComponent memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefRepository {
			updateMemberAccessRefRepo, _ := memberaccessrefdomainrepositories.NewUpdateMemberAccessRefRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return updateMemberAccessRefRepo
		},
	)
}
