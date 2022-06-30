package supplyorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	databasesupplyOrderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositories "github.com/horeekaa/backend/features/supplyOrders/data/repositories"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/utils"
)

type ProposeUpdateSupplyOrderDependency struct{}

func (_ *ProposeUpdateSupplyOrderDependency) Bind() {
	container.Singleton(
		func(
			supplyOrderDataSource databasesupplyOrderdatasourceinterfaces.SupplyOrderDataSource,
			supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
			supplyOrderDataLoader supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader,
		) supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderTransactionComponent {
			proposeUpdateSupplyOrderComponent, _ := supplyorderdomainrepositories.NewProposeUpdateSupplyOrderTransactionComponent(
				supplyOrderDataSource,
				supplyOrderItemDataSource,
				loggingDataSource,
				mapProcessorUtility,
				supplyOrderDataLoader,
			)
			return proposeUpdateSupplyOrderComponent
		},
	)

	container.Transient(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
			supplyOrderDataSource databasesupplyOrderdatasourceinterfaces.SupplyOrderDataSource,
			createPaymentRepository paymentdomainrepositoryinterfaces.CreatePaymentRepository,
			proposeUpdatePaymentRepository paymentdomainrepositoryinterfaces.ProposeUpdatePaymentRepository,
			trxComponent supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderTransactionComponent,
			createSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent,
			proposeUpdateSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemTransactionComponent,
			approveUpdateSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent,
			createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderRepository {
			proposeUpdateSupplyOrderRepo, _ := supplyorderdomainrepositories.NewProposeUpdateSupplyOrderRepository(
				memberAccessDataSource,
				supplyOrderItemDataSource,
				supplyOrderDataSource,
				createPaymentRepository,
				proposeUpdatePaymentRepository,
				trxComponent,
				createSupplyOrderItemComponent,
				proposeUpdateSupplyOrderItemComponent,
				approveUpdateSupplyOrderItemComponent,
				createNotificationComponent,
				mongoDBTransaction,
			)
			return proposeUpdateSupplyOrderRepo
		},
	)
}
