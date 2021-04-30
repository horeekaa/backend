package organizationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositories "github.com/horeekaa/backend/features/organizations/data/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
)

type GetAllOrganizationDependency struct{}

func (_ *GetAllOrganizationDependency) bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
		) organizationdomainrepositoryinterfaces.GetAllOrganizationRepository {
			getAllOrganizationRepo, _ := organizationdomainrepositories.NewGetAllOrganizationRepository(
				organizationDataSource,
			)
			return getAllOrganizationRepo
		},
	)
}
