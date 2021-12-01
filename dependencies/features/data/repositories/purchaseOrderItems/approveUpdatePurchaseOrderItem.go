package purchaseorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositories "github.com/horeekaa/backend/features/purchaseOrderItems/data/repositories"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
)

type ApproveUpdatePurchaseOrderItemDependency struct{}

func (_ *ApproveUpdatePurchaseOrderItemDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			approveUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent {
			approveUpdatePurchaseOrderItemComponent, _ := purchaseorderitemdomainrepositories.NewApproveUpdatePurchaseOrderItemTransactionComponent(
				purchaseOrderItemDataSource,
				loggingDataSource,
				approveUpdateDescriptivePhotoComponent,
				mapProcessorUtility,
			)
			return approveUpdatePurchaseOrderItemComponent
		},
	)
}
