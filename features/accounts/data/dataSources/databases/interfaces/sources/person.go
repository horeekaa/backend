package databaseaccountdatasourceinterfaces

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type PersonDataSource interface {
	GetMongoDataSource() mongodbaccountdatasourceinterfaces.PersonDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbaccountdatasourceinterfaces.PersonDataSourceMongo) bool
}
