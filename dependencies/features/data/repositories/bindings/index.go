package repositoriesdependencies

import (
	accountdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/accounts"
	addressregiongroupdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/addressRegionGroups"
	addressdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/addresses"
	addressdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/addresses/utils"
	descriptivephotodomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/descriptivePhotos"
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
	productvariantdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/productVariants"
	productdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/products"
	purchaseorderitemdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/purchaseOrderItems"
	purchaseorderitemdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/purchaseOrderItems/utils"
	purchaseorderdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/purchaseOrders"
	purchaseorderdomainrepositoryutilitydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/purchaseOrders/utils"
	purchaseordertosupplydomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/purchaseOrdersToSupply"
	taggingdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/taggings"
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
		&addressdomainrepositorydependencies.UpdateAddressDependency{},

		&addressregiongroupdomainrepositorydependencies.CreateAddressRegionGroupDependency{},
		&addressregiongroupdomainrepositorydependencies.ProposeUpdateAddressRegionGroupDependency{},
		&addressregiongroupdomainrepositorydependencies.ApproveUpdateAddressRegionGroupDependency{},
		&addressregiongroupdomainrepositorydependencies.GetAllAddressRegionGroupDependency{},
		&addressregiongroupdomainrepositorydependencies.GetAddressRegionGroupDependency{},

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

		&taggingdomainrepositorydependencies.BulkCreateTaggingDependency{},
		&taggingdomainrepositorydependencies.GetAllTaggingDependency{},
		&taggingdomainrepositorydependencies.GetTaggingDependency{},
		&taggingdomainrepositorydependencies.BulkProposeUpdateTaggingDependency{},
		&taggingdomainrepositorydependencies.BulkApproveUpdateTaggingDependency{},

		&descriptivephotodomainrepositorydependencies.CreateDescriptivePhotoDependency{},
		&descriptivephotodomainrepositorydependencies.UpdateDescriptivePhotoDependency{},
		&descriptivephotodomainrepositorydependencies.GetDescriptivePhotoDependency{},

		&organizationdomainrepositorydependencies.CreateOrganizationDependency{},
		&organizationdomainrepositorydependencies.GetAllOrganizationDependency{},
		&organizationdomainrepositorydependencies.GetOrganizationDependency{},
		&organizationdomainrepositorydependencies.ProposeUpdateOrganizationDependency{},
		&organizationdomainrepositorydependencies.ApproveUpdateOrganizationDependency{},

		&productvariantdomainrepositorydependencies.CreateProductVariantDependency{},
		&productvariantdomainrepositorydependencies.UpdateProductVariantDependency{},
		&productvariantdomainrepositorydependencies.GetProductVariantDependency{},

		&productdomainrepositorydependencies.CreateProductDependency{},
		&productdomainrepositorydependencies.ProposeUpdateProductDependency{},
		&productdomainrepositorydependencies.ApproveUpdateProductDependency{},
		&productdomainrepositorydependencies.GetAllProductDependency{},
		&productdomainrepositorydependencies.GetProductDependency{},

		&purchaseorderitemdomainrepositoryutilitydependencies.PurchaseOrderItemLoaderDependency{},
		&purchaseorderitemdomainrepositorydependencies.CreatePurchaseOrderItemDependency{},
		&purchaseorderitemdomainrepositorydependencies.UpdatePurchaseOrderItemDependency{},
		&purchaseorderitemdomainrepositorydependencies.GetPurchaseOrderItemDependency{},

		&purchaseorderdomainrepositoryutilitydependencies.PurchaseOrderLoaderDependency{},
		&purchaseorderdomainrepositorydependencies.CreatePurchaseOrderDependency{},
		&purchaseorderdomainrepositorydependencies.GetAllPurchaseOrderDependency{},
		&purchaseorderdomainrepositorydependencies.GetPurchaseOrderDependency{},
		&purchaseorderdomainrepositorydependencies.ApproveUpdatePurchaseOrderDependency{},
		&purchaseorderdomainrepositorydependencies.ProposeUpdatePurchaseOrderDependency{},

		&purchaseordertosupplydomainrepositorydependencies.CreatePurchaseOrderToSupplyDependency{},
		&purchaseordertosupplydomainrepositorydependencies.GetAllPurchaseOrderToSupplyDependency{},
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

		&moudomainrepositoryutilitydependencies.PartyLoaderDependency{},
		&moudomainrepositorydependencies.CreateMouDependency{},
		&moudomainrepositorydependencies.ProposeUpdateMouDependency{},
		&moudomainrepositorydependencies.ApproveUpdateMouDependency{},
		&moudomainrepositorydependencies.GetAllMouDependency{},
		&moudomainrepositorydependencies.GetMouDependency{},

		&mouitemdomainrepositoryutilitydependencies.AgreedProductLoaderDependency{},
		&mouitemdomainrepositorydependencies.CreateMouItemDependency{},
		&mouitemdomainrepositorydependencies.UpdateMouItemDependency{},
		&mouitemdomainrepositorydependencies.GetMouItemDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
