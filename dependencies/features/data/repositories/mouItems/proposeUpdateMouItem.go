package mouitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositories "github.com/horeekaa/backend/features/mouItems/data/repositories"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories/utils"
)

type ProposeUpdateMouItemDependency struct{}

func (_ *ProposeUpdateMouItemDependency) Bind() {
	container.Singleton(
		func(
			mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			agreedProductLoader mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) mouitemdomainrepositoryinterfaces.ProposeUpdateMouItemTransactionComponent {
			updateMouItemComponent, _ := mouitemdomainrepositories.NewProposeUpdateMouItemTransactionComponent(
				mouItemDataSource,
				loggingDataSource,
				agreedProductLoader,
				mapProcessorUtility,
			)
			return updateMouItemComponent
		},
	)
}
