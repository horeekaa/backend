package purchaseorderdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositories "github.com/horeekaa/backend/features/purchaseOrders/data/repositories"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
)

type UpdatePurchaseOrderByCronDependency struct{}

func (_ *UpdatePurchaseOrderByCronDependency) Bind() {
	container.Transient(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			mouDataSource databasemoudatasourceinterfaces.MouDataSource,
			purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
			purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
			createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
		) purchaseorderdomainrepositoryinterfaces.UpdatePurchaseOrderByCronRepository {
			updatePurchaseOrderByCronRepo, _ := purchaseorderdomainrepositories.NewUpdatePurchaseOrderByCronRepository(
				memberAccessDataSource,
				mouDataSource,
				purchaseOrderDataSource,
				purchaseOrderItemDataSource,
				createNotificationComponent,
			)
			return updatePurchaseOrderByCronRepo
		},
	)
}
