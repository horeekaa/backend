package supplyorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositories "github.com/horeekaa/backend/features/supplyOrders/data/repositories"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/utils"
)

type ApproveUpdateSupplyOrderDependency struct{}

func (_ *ApproveUpdateSupplyOrderDependency) Bind() {
	container.Singleton(
		func(
			supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
			supplyOrderDataLoader supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader,
		) supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderTransactionComponent {
			approveUpdateSupplyOrderComponent, _ := supplyorderdomainrepositories.NewApproveUpdateSupplyOrderTransactionComponent(
				supplyOrderDataSource,
				loggingDataSource,
				mapProcessorUtility,
				supplyOrderDataLoader,
			)
			return approveUpdateSupplyOrderComponent
		},
	)

	container.Transient(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
			approveUpdatePaymentComponent paymentdomainrepositoryinterfaces.ApproveUpdatePaymentTransactionComponent,
			approveUpdatesupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent,
			trxComponent supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderTransactionComponent,
			createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderRepository {
			approveUpdatesupplyOrderRepo, _ := supplyorderdomainrepositories.NewApproveUpdateSupplyOrderRepository(
				memberAccessDataSource,
				supplyOrderDataSource,
				approveUpdatePaymentComponent,
				approveUpdatesupplyOrderItemComponent,
				trxComponent,
				createNotificationComponent,
				mongoDBTransaction,
			)
			return approveUpdatesupplyOrderRepo
		},
	)
}
