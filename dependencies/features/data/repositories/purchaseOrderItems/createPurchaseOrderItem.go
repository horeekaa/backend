package purchaseorderitemdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositories "github.com/horeekaa/backend/features/purchaseOrderItems/data/repositories"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/utils"
)

type CreatePurchaseOrderItemDependency struct{}

func (_ *CreatePurchaseOrderItemDependency) Bind() {
	container.Singleton(
		func(
			purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
			purchaseOrderItemLoader purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader,
		) purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent {
			createPurchaseOrderItemComponent, _ := purchaseorderitemdomainrepositories.NewCreatePurchaseOrderItemTransactionComponent(
				purchaseOrderItemDataSource,
				loggingDataSource,
				createDescriptivePhotoComponent,
				purchaseOrderItemLoader,
			)
			return createPurchaseOrderItemComponent
		},
	)
}
