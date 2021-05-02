package mongodborganizationdatasourcedependencies

import (
	container "github.com/golobby/container/v2/"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	mongodborganizationdatasources "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/mongodb"
	mongodborganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/mongodb/interfaces"
	databaseorganizationdatasources "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/sources"
)

type OrganizationDataSourceDependency struct{}

func (orgDataSourceDpdcy *OrganizationDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodborganizationdatasourceinterfaces.OrganizationDataSourceMongo {
			organizationDataSourceMongo, _ := mongodborganizationdatasources.NewOrganizationDataSourceMongo(basicOperation)
			return organizationDataSourceMongo
		},
	)

	container.Singleton(
		func(organizationDataSourceMongo mongodborganizationdatasourceinterfaces.OrganizationDataSourceMongo) databaseorganizationdatasourceinterfaces.OrganizationDataSource {
			organizationRepo, _ := databaseorganizationdatasources.NewOrganizationDataSource()
			organizationRepo.SetMongoDataSource(organizationDataSourceMongo)
			return organizationRepo
		},
	)
}
