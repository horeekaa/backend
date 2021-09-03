package moudomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	moudomainrepositories "github.com/horeekaa/backend/features/mous/data/repositories"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moudomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories/utils"
)

type CreateMouDependency struct{}

func (_ *CreateMouDependency) Bind() {
	container.Singleton(
		func(
			mouDataSource databasemoudatasourceinterfaces.MouDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
			partyLoader moudomainrepositoryutilityinterfaces.PartyLoader,
		) moudomainrepositoryinterfaces.CreateMouTransactionComponent {
			createmouComponent, _ := moudomainrepositories.NewCreateMouTransactionComponent(
				mouDataSource,
				loggingDataSource,
				structFieldIteratorUtility,
				partyLoader,
			)
			return createmouComponent
		},
	)

	container.Transient(
		func(
			trxComponent moudomainrepositoryinterfaces.CreateMouTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) moudomainrepositoryinterfaces.CreateMouRepository {
			createMouRepo, _ := moudomainrepositories.NewCreateMouRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return createMouRepo
		},
	)
}
