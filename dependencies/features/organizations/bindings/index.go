package organizationdependencies

import (
	mongodborganizationdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/organizations/data/dataSources/databases/mongodb"
	organizationdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/organizations/data/repositories"
	organizationpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/organizations/domain/usecases"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type OrganizationDependency struct{}

func (_ *OrganizationDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&mongodborganizationdatasourcedependencies.OrganizationDataSourceDependency{},

		&organizationdomainrepositorydependencies.CreateOrganizationDependency{},
		&organizationdomainrepositorydependencies.GetAllOrganizationDependency{},
		&organizationdomainrepositorydependencies.GetOrganizationDependency{},
		&organizationdomainrepositorydependencies.UpdateOrganizationDependency{},

		&organizationpresentationusecasedependencies.CreateOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.GetAllOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.GetOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.UpdateOrganizationUsecaseDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
