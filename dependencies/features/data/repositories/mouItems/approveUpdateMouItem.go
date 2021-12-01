package mouitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositories "github.com/horeekaa/backend/features/mouItems/data/repositories"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
)

type ApproveUpdateMouItemDependency struct{}

func (_ *ApproveUpdateMouItemDependency) Bind() {
	container.Singleton(
		func(
			mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) mouitemdomainrepositoryinterfaces.ApproveUpdateMouItemTransactionComponent {
			approveUpdateMouItemComponent, _ := mouitemdomainrepositories.NewApproveUpdateMouItemTransactionComponent(
				mouItemDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdateMouItemComponent
		},
	)
}
