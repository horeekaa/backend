package databaseaccountdatasourceinterfaces

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type AccountDataSource interface {
	GetMongoDataSource() mongodbaccountdatasourceinterfaces.AccountDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbaccountdatasourceinterfaces.AccountDataSourceMongo) bool
}
