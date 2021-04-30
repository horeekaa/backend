package databaseorganizationdatasources

import (
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	mongodborganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/mongodb/interfaces"
)

type organizationDataSource struct {
	organizationDataSourceRepoMongo mongodborganizationdatasourceinterfaces.OrganizationDataSourceMongo
}

func (orgDataSource *organizationDataSource) SetMongoDataSource(mongoDataSource mongodborganizationdatasourceinterfaces.OrganizationDataSourceMongo) bool {
	orgDataSource.organizationDataSourceRepoMongo = mongoDataSource
	return true
}

func (orgDataSource *organizationDataSource) GetMongoDataSource() mongodborganizationdatasourceinterfaces.OrganizationDataSourceMongo {
	return orgDataSource.organizationDataSourceRepoMongo
}

func NewOrganizationDataSource() (databaseorganizationdatasourceinterfaces.OrganizationDataSource, error) {
	return &organizationDataSource{}, nil
}
