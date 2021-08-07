package repositoriesdependencies

import (
	accountdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/accounts"
	descriptivephotodomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/descriptivePhotos"
	loggingdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/loggings"
	memberaccessrefdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/memberAccessRefs"
	memberaccessdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/memberAccesses"
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

		&organizationdomainrepositorydependencies.CreateOrganizationDependency{},
		&organizationdomainrepositorydependencies.GetAllOrganizationDependency{},
		&organizationdomainrepositorydependencies.GetOrganizationDependency{},
		&organizationdomainrepositorydependencies.ProposeUpdateOrganizationDependency{},
		&organizationdomainrepositorydependencies.ApproveUpdateOrganizationDependency{},

		&descriptivephotodomainrepositorydependencies.CreateDescriptivePhotoDependency{},
		&descriptivephotodomainrepositorydependencies.UpdateDescriptivePhotoDependency{},
		&descriptivephotodomainrepositorydependencies.GetDescriptivePhotoDependency{},

		&productdomainrepositorydependencies.CreateProductDependency{},
		&productdomainrepositorydependencies.ProposeUpdateProductDependency{},
		&productdomainrepositorydependencies.ApproveUpdateProductDependency{},
		&productdomainrepositorydependencies.GetAllProductDependency{},
		&productdomainrepositorydependencies.GetProductDependency{},

		&productvariantdomainrepositorydependencies.CreateProductVariantDependency{},
		&productvariantdomainrepositorydependencies.UpdateProductVariantDependency{},
		&productvariantdomainrepositorydependencies.GetProductVariantDependency{},

		&tagdomainrepositorydependencies.CreateTagDependency{},
		&tagdomainrepositorydependencies.ProposeUpdateTagDependency{},
		&tagdomainrepositorydependencies.ApproveUpdateTagDependency{},
		&tagdomainrepositorydependencies.GetAllTagDependency{},
		&tagdomainrepositorydependencies.GetTagDependency{},

		&taggingdomainrepositorydependencies.BulkCreateTaggingDependency{},
		&taggingdomainrepositorydependencies.GetAllTaggingDependency{},
		&taggingdomainrepositorydependencies.GetTaggingDependency{},
		&taggingdomainrepositorydependencies.BulkProposeUpdateTaggingDependency{},
		&taggingdomainrepositorydependencies.BulkApproveUpdateTaggingDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
