package supplyorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databasesupplyorderItemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositories "github.com/horeekaa/backend/features/supplyOrderItems/data/repositories"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
)

type ProposeUpdateSupplyOrderItemPickUpDependency struct{}

func (_ *ProposeUpdateSupplyOrderItemPickUpDependency) Bind() {
	container.Transient(
		func(
			supplyOrderItemDataSource databasesupplyorderItemdatasourceinterfaces.SupplyOrderItemDataSource,
			proposeUpdatesupplyOrderItemTransactionComponent supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemPickUpRepository {
			proposeUpdateSupplyOrderItemPickUpRepo, _ := supplyorderitemdomainrepositories.NewProposeUpdateSupplyOrderItemPickUpRepository(
				supplyOrderItemDataSource,
				proposeUpdatesupplyOrderItemTransactionComponent,
				mongoDBTransaction,
			)
			return proposeUpdateSupplyOrderItemPickUpRepo
		},
	)
}
