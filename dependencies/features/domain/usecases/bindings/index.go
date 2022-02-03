package usecasesdependencies

import (
	accountpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/accounts"
	addressregiongrouppresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/addressRegionGroups"
	addresspresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/addresses"
	descriptivephotopresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/descriptivePhotos"
	invoicepresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/invoices"
	loggingpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/loggings"
	memberaccessrefpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/memberAccessRefs"
	memberaccesspresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/memberAccesses"
	mouitempresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/mouItems"
	moupresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/mous"
	notificationpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/notifications"
	organizationpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/organizations"
	productvariantpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/productVariants"
	productpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/products"
	purchaseorderitempresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/purchaseOrderItems"
	purchaseorderpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/purchaseOrders"
	purchaseordertosupplypresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/purchaseOrdersToSupply"
	supplyorderitempresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/supplyOrderItems"
	supplyorderpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/supplyOrders"
	taggingpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/taggings"
	tagpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/tags"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type UsecasesDependency struct{}

func (_ *UsecasesDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&accountpresentationusecasedependencies.GetAccountUsecaseDependency{},
		&accountpresentationusecasedependencies.GetPersonDataFromAccountUsecaseDependency{},
		&accountpresentationusecasedependencies.LoginUsecaseDependency{},
		&accountpresentationusecasedependencies.LogoutUsecaseDependency{},
		&accountpresentationusecasedependencies.GetAuthUserAndAttachToCtxUsecaseDependency{},

		&addresspresentationusecasedependencies.GetAddressUsecaseDependency{},

		&addressregiongrouppresentationusecasedependencies.CreateAddressRegionGroupUsecaseDependency{},
		&addressregiongrouppresentationusecasedependencies.GetAllAddressRegionGroupUsecaseDependency{},
		&addressregiongrouppresentationusecasedependencies.GetAddressRegionGroupUsecaseDependency{},
		&addressregiongrouppresentationusecasedependencies.UpdateAddressRegionGroupUsecaseDependency{},

		&invoicepresentationusecasedependencies.CreateInvoiceUsecaseDependency{},
		&invoicepresentationusecasedependencies.GetAllInvoiceUsecaseDependency{},
		&invoicepresentationusecasedependencies.GetInvoiceUsecaseDependency{},
		&invoicepresentationusecasedependencies.UpdateInvoiceUsecaseDependency{},

		&loggingpresentationusecasedependencies.GetLoggingUsecaseDependency{},

		&memberaccesspresentationusecasedependencies.CreateMemberAccessUsecaseDependency{},
		&memberaccesspresentationusecasedependencies.GetAllMemberAccessUsecaseDependency{},
		&memberaccesspresentationusecasedependencies.GetMemberAccessUsecaseDependency{},
		&memberaccesspresentationusecasedependencies.UpdateMemberAccessUsecaseDependency{},

		&memberaccessrefpresentationusecasedependencies.CreateMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.GetAllMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.GetMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.UpdateMemberAccessRefUsecaseDependency{},

		&organizationpresentationusecasedependencies.CreateOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.GetAllOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.GetOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.UpdateOrganizationUsecaseDependency{},

		&descriptivephotopresentationusecasedependencies.GetDescriptivePhotoUsecaseDependency{},

		&productpresentationusecasedependencies.CreateProductUsecaseDependency{},
		&productpresentationusecasedependencies.GetAllProductUsecaseDependency{},
		&productpresentationusecasedependencies.GetProductUsecaseDependency{},
		&productpresentationusecasedependencies.UpdateProductUsecaseDependency{},

		&productvariantpresentationusecasedependencies.GetProductVariantUsecaseDependency{},

		&tagpresentationusecasedependencies.CreateTagUsecaseDependency{},
		&tagpresentationusecasedependencies.GetAllTagUsecaseDependency{},
		&tagpresentationusecasedependencies.GetTagUsecaseDependency{},
		&tagpresentationusecasedependencies.UpdateTagUsecaseDependency{},

		&taggingpresentationusecasedependencies.BulkCreateTaggingUsecaseDependency{},
		&taggingpresentationusecasedependencies.BulkUpdateTaggingUsecaseDependency{},
		&taggingpresentationusecasedependencies.GetAllTaggingUsecaseDependency{},
		&taggingpresentationusecasedependencies.GetTaggingUsecaseDependency{},

		&notificationpresentationusecasedependencies.GetAllNotificationUsecaseDependency{},
		&notificationpresentationusecasedependencies.BulkUpdateNotificationUsecaseDependency{},

		&moupresentationusecasedependencies.CreateMouUsecaseDependency{},
		&moupresentationusecasedependencies.GetAllMouUsecaseDependency{},
		&moupresentationusecasedependencies.GetMouUsecaseDependency{},
		&moupresentationusecasedependencies.UpdateMouUsecaseDependency{},

		&mouitempresentationusecasedependencies.GetMouItemUsecaseDependency{},

		&purchaseorderpresentationusecasedependencies.CreatePurchaseOrderUsecaseDependency{},
		&purchaseorderpresentationusecasedependencies.GetAllPurchaseOrderUsecaseDependency{},
		&purchaseorderpresentationusecasedependencies.GetPurchaseOrderUsecaseDependency{},
		&purchaseorderpresentationusecasedependencies.UpdatePurchaseOrderUsecaseDependency{},

		&purchaseorderitempresentationusecasedependencies.UpdatePurchaseOrderItemDeliveryUsecaseDependency{},
		&purchaseorderitempresentationusecasedependencies.GetPurchaseOrderItemUsecaseDependency{},

		&purchaseordertosupplypresentationusecasedependencies.GetAllPurchaseOrderToSupplyUsecaseDependency{},
		&purchaseordertosupplypresentationusecasedependencies.GetPurchaseOrderToSupplyUsecaseDependency{},
		&purchaseordertosupplypresentationusecasedependencies.ProcessPurchaseOrderToSupplyUsecaseDependency{},

		&supplyorderpresentationusecasedependencies.CreateSupplyOrderUsecaseDependency{},
		&supplyorderpresentationusecasedependencies.GetSupplyOrderUsecaseDependency{},
		&supplyorderpresentationusecasedependencies.GetAllSupplyOrderUsecaseDependency{},
		&supplyorderpresentationusecasedependencies.UpdateSupplyOrderUsecaseDependency{},

		&supplyorderitempresentationusecasedependencies.UpdateSupplyOrderItemPickUpUsecaseDependency{},
		&supplyorderitempresentationusecasedependencies.GetSupplyOrderItemUsecaseDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
