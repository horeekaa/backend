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

type ProposeUpdateMouDependency struct{}

func (_ *ProposeUpdateMouDependency) Bind() {
	container.Singleton(
		func(
			mouDataSource databasemoudatasourceinterfaces.MouDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
			structComparisonUtility coreutilityinterfaces.StructComparisonUtility,
			partyLoader moudomainrepositoryutilityinterfaces.PartyLoader,
		) moudomainrepositoryinterfaces.ProposeUpdateMouTransactionComponent {
			proposeUpdateMouComponent, _ := moudomainrepositories.NewProposeUpdateMouTransactionComponent(
				mouDataSource,
				loggingDataSource,
				mapProcessorUtility,
				structComparisonUtility,
				partyLoader,
			)
			return proposeUpdateMouComponent
		},
	)

	container.Transient(
		func(
			trxComponent moudomainrepositoryinterfaces.ProposeUpdateMouTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) moudomainrepositoryinterfaces.ProposeUpdateMouRepository {
			proposeUpdateMouRepo, _ := moudomainrepositories.NewProposeUpdateMouRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return proposeUpdateMouRepo
		},
	)
}
