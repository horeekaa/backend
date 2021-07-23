package repositoriesdependencies

import (
	accountdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/accounts"
	descriptivephotodomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/descriptivePhotos"
	loggingdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/loggings"
	memberaccessrefdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/memberAccessRefs"
	memberaccessdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/memberAccesses"
	organizationdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/organizations"
	productdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/products"
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
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
