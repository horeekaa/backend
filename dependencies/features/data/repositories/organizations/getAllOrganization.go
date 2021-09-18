package organizationdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositories "github.com/horeekaa/backend/features/organizations/data/repositories"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
)

type GetAllOrganizationDependency struct{}

func (_ *GetAllOrganizationDependency) Bind() {
	container.Singleton(
		func(
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) organizationdomainrepositoryinterfaces.GetAllOrganizationRepository {
			getAllOrganizationRepo, _ := organizationdomainrepositories.NewGetAllOrganizationRepository(
				organizationDataSource,
				mongoQueryBuilder,
			)
			return getAllOrganizationRepo
		},
	)
}
