package memberaccessrefdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositories "github.com/horeekaa/backend/features/memberAccessRefs/data/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
)

type CreateMemberAccessRefDependency struct{}

func (_ *CreateMemberAccessRefDependency) Bind() {
	container.Singleton(
		func(
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
		) memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefTransactionComponent {
			createMemberAccessRefComponent, _ := memberaccessrefdomainrepositories.NewCreateMemberAccessRefTransactionComponent(
				memberAccessRefDataSource,
				loggingDataSource,
			)
			return createMemberAccessRefComponent
		},
	)

	container.Transient(
		func(
			trxComponent memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefRepository {
			updateMemberAccessRefRepo, _ := memberaccessrefdomainrepositories.NewCreateMemberAccessRefRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return updateMemberAccessRefRepo
		},
	)
}
