package mouitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositories "github.com/horeekaa/backend/features/mouItems/data/repositories"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories/utils"
)

type UpdateMouItemDependency struct{}

func (_ *UpdateMouItemDependency) Bind() {
	container.Singleton(
		func(
			mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
			agreedProductLoader mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader,
		) mouitemdomainrepositoryinterfaces.UpdateMouItemTransactionComponent {
			updateMouItemComponent, _ := mouitemdomainrepositories.NewUpdateMouItemTransactionComponent(
				mouItemDataSource,
				agreedProductLoader,
			)
			return updateMouItemComponent
		},
	)
}
