package memberaccessrefdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/transaction"
	memberaccessrefdomainrepositories "github.com/horeekaa/backend/features/memberAccessRefs/data/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberaccessrefs/data/dataSources/databases/interfaces/sources"
)

type UpdateMemberAccessRefDependency struct{}

func (_ UpdateMemberAccessRefDependency) bind() {
	container.Singleton(
		func(
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
		) memberaccessrefdomainrepositoryinterfaces.UpdateMemberAccessRefTransactionComponent {
			updateMemberAccessRefComponent, _ := memberaccessrefdomainrepositories.NewUpdateMemberAccessRefTransactionComponent(
				memberAccessRefDataSource,
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
