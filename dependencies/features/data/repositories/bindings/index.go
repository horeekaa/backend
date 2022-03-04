package repositoriesdependencies

import (
	accountdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/accounts"
	addressregiongroupdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/addressRegionGroups"
	addressdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/addresses"
	addressdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/addresses/utils"
	descriptivephotodomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/descriptivePhotos"
	invoicedomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/invoices"
	loggingdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/loggings"
	memberaccessrefdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/memberAccessRefs"
	memberaccessdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/memberAccesses"
	mouitemdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/mouItems"
	mouitemdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/mouItems/utils"
	moudomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/mous"
	moudomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/mous/utils"
	notificationdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/notifications"
	notificationdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/notifications/utils"
	notificationdomainrepositoryloaderutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/notifications/utils/payloadLoaders"
	organizationdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/organizations"
	paymentdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/payments"
	paymentdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/payments/utils"
	productvariantdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/productVariants"
	productdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/products"
	purchaseorderitemdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/purchaseOrderItems"
	purchaseorderitemdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/purchaseOrderItems/utils"
	purchaseorderdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/purchaseOrders"
	purchaseorderdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/purchaseOrders/utils"
	purchaseordertosupplydomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/purchaseOrdersToSupply"
	supplyorderitemdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/supplyOrderItems"
	supplyorderitemdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/supplyOrderItems/utils"
	supplyorderdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/supplyOrders"
	supplyorderdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/supplyOrders/utils"
	taggingdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/taggings"
	taggingdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/taggings/utils"
	tagdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/tags"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type RepositoriesDependency struct{}

func (_ *RepositoriesDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&accountdomainrepositorydependencies.GetAccountDependency{},
		&accountdomainrepositorydependencies.GetPersonDataFromAccountDependency{},
		&accountdomainrepositorydependencies.CreateAccountFromAuthDataDependency{},
		&accountdomainrepositorydependencies.GetAccountFromAuthDataDependency{},
		&accountdomainrepositorydependencies.GetUserFromAuthHeaderDependency{},
		&accountdomainrepositorydependencies.ManageAccountDeviceTokenDependency{},

		&addressdomainrepositoryutilitydependencies.AddressLoaderDependency{},
		&addressdomainrepositorydependencies.CreateAddressDependency{},
		&addressdomainrepositorydependencies.GetAddressDependency{},
		&addressdomainrepositorydependencies.ProposeUpdateAddressDependency{},
		&addressdomainrepositorydependencies.ApproveUpdateAddressDependency{},

		&addressregiongroupdomainrepositorydependencies.CreateAddressRegionGroupDependency{},
		&addressregiongroupdomainrepositorydependencies.ProposeUpdateAddressRegionGroupDependency{},
		&addressregiongroupdomainrepositorydependencies.ApproveUpdateAddressRegionGroupDependency{},
		&addressregiongroupdomainrepositorydependencies.GetAllAddressRegionGroupDependency{},
		&addressregiongroupdomainrepositorydependencies.GetAddressRegionGroupDependency{},

		&invoicedomainrepositorydependencies.CreateInvoiceDependency{},
		&invoicedomainrepositorydependencies.GetAllInvoiceDependency{},
		&invoicedomainrepositorydependencies.GetInvoiceDependency{},
		&invoicedomainrepositorydependencies.UpdateInvoiceDependency{},

		&loggingdomainrepositorydependencies.GetLoggingDependency{},

		&memberaccessdomainrepositorydependencies.CreateMemberAccessDependency{},
		&memberaccessdomainrepositorydependencies.GetAccountMemberAccessDependency{},
		&memberaccessdomainrepositorydependencies.GetAllMemberAccessDependency{},
		&memberaccessdomainrepositorydependencies.ProposeUpdateMemberAccessDependency{},
		&memberaccessdomainrepositorydependencies.ApproveUpdateMemberAccessDependency{},

		&memberaccessrefdomainrepositorydependencies.CreateMemberAccessRefDependency{},
		&memberaccessrefdomainrepositorydependencies.GetAllMemberAccessRefDependency{},
		&memberaccessrefdomainrepositorydependencies.GetMemberAccessRefDependency{},
		&memberaccessrefdomainrepositorydependencies.ProposeUpdateMemberAccessRefDependency{},
		&memberaccessrefdomainrepositorydependencies.ApproveUpdateMemberAccessRefDependency{},

		&taggingdomainrepositoryutilitydependencies.TaggingLoaderDependency{},
		&taggingdomainrepositorydependencies.BulkCreateTaggingDependency{},
		&taggingdomainrepositorydependencies.GetAllTaggingDependency{},
		&taggingdomainrepositorydependencies.GetTaggingDependency{},
		&taggingdomainrepositorydependencies.BulkProposeUpdateTaggingDependency{},
		&taggingdomainrepositorydependencies.BulkApproveUpdateTaggingDependency{},

		&descriptivephotodomainrepositorydependencies.CreateDescriptivePhotoDependency{},
		&descriptivephotodomainrepositorydependencies.ProposeUpdateDescriptivePhotoDependency{},
		&descriptivephotodomainrepositorydependencies.ApproveUpdateDescriptivePhotoDependency{},
		&descriptivephotodomainrepositorydependencies.GetDescriptivePhotoDependency{},

		&organizationdomainrepositorydependencies.CreateOrganizationDependency{},
		&organizationdomainrepositorydependencies.GetAllOrganizationDependency{},
		&organizationdomainrepositorydependencies.GetOrganizationDependency{},
		&organizationdomainrepositorydependencies.ProposeUpdateOrganizationDependency{},
		&organizationdomainrepositorydependencies.ApproveUpdateOrganizationDependency{},

		&paymentdomainrepositoryutilitydependencies.PaymentLoaderDependency{},
		&paymentdomainrepositorydependencies.CreatePaymentDependency{},
		&paymentdomainrepositorydependencies.GetAllPaymentDependency{},
		&paymentdomainrepositorydependencies.GetPaymentDependency{},
		&paymentdomainrepositorydependencies.ApproveUpdatePaymentDependency{},
		&paymentdomainrepositorydependencies.ProposeUpdatePaymentDependency{},

		&productvariantdomainrepositorydependencies.CreateProductVariantDependency{},
		&productvariantdomainrepositorydependencies.ApproveUpdateProductVariantDependency{},
		&productvariantdomainrepositorydependencies.ProposeUpdateProductVariantDependency{},
		&productvariantdomainrepositorydependencies.GetProductVariantDependency{},

		&productdomainrepositorydependencies.CreateProductDependency{},
		&productdomainrepositorydependencies.ProposeUpdateProductDependency{},
		&productdomainrepositorydependencies.ApproveUpdateProductDependency{},
		&productdomainrepositorydependencies.GetAllProductDependency{},
		&productdomainrepositorydependencies.GetProductDependency{},

		&purchaseorderitemdomainrepositoryutilitydependencies.PurchaseOrderItemLoaderDependency{},
		&purchaseorderitemdomainrepositorydependencies.CreatePurchaseOrderItemDependency{},
		&purchaseorderitemdomainrepositorydependencies.ApproveUpdatePurchaseOrderItemDependency{},
		&purchaseorderitemdomainrepositorydependencies.ProposeUpdatePurchaseOrderItemDependency{},
		&purchaseorderitemdomainrepositorydependencies.GetPurchaseOrderItemDependency{},
		&purchaseorderitemdomainrepositorydependencies.GetAllPurchaseOrderItemDependency{},
		&purchaseorderitemdomainrepositorydependencies.ProposeUpdatePurchaseOrderItemDeliveryDependency{},

		&purchaseorderdomainrepositoryutilitydependencies.PurchaseOrderLoaderDependency{},
		&purchaseorderdomainrepositorydependencies.CreatePurchaseOrderDependency{},
		&purchaseorderdomainrepositorydependencies.GetAllPurchaseOrderDependency{},
		&purchaseorderdomainrepositorydependencies.GetPurchaseOrderDependency{},
		&purchaseorderdomainrepositorydependencies.ApproveUpdatePurchaseOrderDependency{},
		&purchaseorderdomainrepositorydependencies.ProposeUpdatePurchaseOrderDependency{},

		&purchaseordertosupplydomainrepositorydependencies.GetAllPurchaseOrderToSupplyDependency{},
		&purchaseordertosupplydomainrepositorydependencies.GetPurchaseOrderToSupplyDependency{},
		&purchaseordertosupplydomainrepositorydependencies.ProcessPurchaseOrderToSupplyDependency{},

		&tagdomainrepositorydependencies.CreateTagDependency{},
		&tagdomainrepositorydependencies.ProposeUpdateTagDependency{},
		&tagdomainrepositorydependencies.ApproveUpdateTagDependency{},
		&tagdomainrepositorydependencies.GetAllTagDependency{},
		&tagdomainrepositorydependencies.GetTagDependency{},

		&notificationdomainrepositoryloaderutilitydependencies.InvitationPayloadLoaderDependency{},
		&notificationdomainrepositoryutilitydependencies.MasterPayloadLoaderDependency{},
		&notificationdomainrepositoryutilitydependencies.NotificationLocalizationBuilderDependency{},
		&notificationdomainrepositorydependencies.CreateNotificationDependency{},
		&notificationdomainrepositorydependencies.BulkUpdateNotificationDependency{},
		&notificationdomainrepositorydependencies.GetAllNotificationDependency{},

		&mouitemdomainrepositoryutilitydependencies.AgreedProductLoaderDependency{},
		&mouitemdomainrepositorydependencies.CreateMouItemDependency{},
		&mouitemdomainrepositorydependencies.ApproveUpdateMouItemDependency{},
		&mouitemdomainrepositorydependencies.ProposeUpdateMouItemDependency{},
		&mouitemdomainrepositorydependencies.GetMouItemDependency{},

		&moudomainrepositoryutilitydependencies.PartyLoaderDependency{},
		&moudomainrepositorydependencies.CreateMouDependency{},
		&moudomainrepositorydependencies.ProposeUpdateMouDependency{},
		&moudomainrepositorydependencies.ApproveUpdateMouDependency{},
		&moudomainrepositorydependencies.GetAllMouDependency{},
		&moudomainrepositorydependencies.GetMouDependency{},

		&supplyorderitemdomainrepositoryutilitydependencies.SupplyOrderItemLoaderDependency{},
		&supplyorderitemdomainrepositorydependencies.CreateSupplyOrderItemDependency{},
		&supplyorderitemdomainrepositorydependencies.ApproveUpdateSupplyOrderItemDependency{},
		&supplyorderitemdomainrepositorydependencies.ProposeUpdateSupplyOrderItemDependency{},
		&supplyorderitemdomainrepositorydependencies.GetSupplyOrderItemDependency{},
		&supplyorderitemdomainrepositorydependencies.ProposeUpdateSupplyOrderItemPickUpDependency{},

		&supplyorderdomainrepositoryutilitydependencies.SupplyOrderLoaderDependency{},
		&supplyorderdomainrepositorydependencies.CreateSupplyOrderDependency{},
		&supplyorderdomainrepositorydependencies.GetSupplyOrderDependency{},
		&supplyorderdomainrepositorydependencies.GetAllSupplyOrderDependency{},
		&supplyorderdomainrepositorydependencies.ApproveUpdateSupplyOrderDependency{},
		&supplyorderdomainrepositorydependencies.ProposeUpdateSupplyOrderDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
