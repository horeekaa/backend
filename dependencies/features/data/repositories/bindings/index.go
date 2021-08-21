package repositoriesdependencies

import (
	accountdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/accounts"
	descriptivephotodomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/descriptivePhotos"
	loggingdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/loggings"
	memberaccessrefdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/memberAccessRefs"
	memberaccessdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/memberAccesses"
	notificationdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/notifications"
	organizationdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/organizations"
	productvariantdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/productVariants"
	productdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/products"
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

		&tagdomainrepositorydependencies.CreateTagDependency{},
		&tagdomainrepositorydependencies.ProposeUpdateTagDependency{},
		&tagdomainrepositorydependencies.ApproveUpdateTagDependency{},
		&tagdomainrepositorydependencies.GetAllTagDependency{},
		&tagdomainrepositorydependencies.GetTagDependency{},

		&notificationdomainrepositorydependencies.NotificationLocalizationBuilderDependency{},
		&notificationdomainrepositorydependencies.CreateNotificationDependency{},
		&notificationdomainrepositorydependencies.BulkUpdateNotificationDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
