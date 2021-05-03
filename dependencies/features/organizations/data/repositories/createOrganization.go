package organizationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositories "github.com/horeekaa/backend/features/organizations/data/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
)

type CreateOrganizationDependency struct{}

func (_ *CreateOrganizationDependency) Bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
		) organizationdomainrepositoryinterfaces.CreateOrganizationRepository {
			createOrganizationRepo, _ := organizationdomainrepositories.NewCreateOrganizationRepository(
				organizationDataSource,
			)
			return createOrganizationRepo
		},
	)
}
