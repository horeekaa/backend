package databaseorganizationdatasourceinterfaces

import (
	mongodborganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/mongodb/interfaces"
)

type OrganizationDataSource interface {
	GetMongoDataSource() mongodborganizationdatasourceinterfaces.OrganizationDataSourceMongo
	SetMongoDataSource(mongoRepo mongodborganizationdatasourceinterfaces.OrganizationDataSourceMongo) bool
}
