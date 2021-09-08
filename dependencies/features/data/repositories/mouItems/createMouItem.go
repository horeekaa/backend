package mouitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasemouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/interfaces/sources"
	mouitemdomainrepositories "github.com/horeekaa/backend/features/mouItems/data/repositories"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories/utils"
)

type CreateMouItemDependency struct{}

func (_ *CreateMouItemDependency) Bind() {
	container.Singleton(
		func(
			mouItemDataSource databasemouitemdatasourceinterfaces.MouItemDataSource,
			agreedProductLoader mouitemdomainrepositoryutilityinterfaces.AgreedProductLoader,
		) mouitemdomainrepositoryinterfaces.CreateMouItemTransactionComponent {
			createMouItemComponent, _ := mouitemdomainrepositories.NewCreateMouItemTransactionComponent(
				mouItemDataSource,
				agreedProductLoader,
			)
			return createMouItemComponent
		},
	)
}
