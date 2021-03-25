package databaseaccountdatasourceinterfaces

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type MemberAccessRefDataSource interface {
	GetMongoDataSource() mongodbaccountdatasourceinterfaces.MemberAccessRefDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbaccountdatasourceinterfaces.MemberAccessRefDataSourceMongo) bool
}
