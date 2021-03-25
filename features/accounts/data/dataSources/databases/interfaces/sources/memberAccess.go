package databaseaccountdatasourceinterfaces

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type MemberAccessDataSource interface {
	GetMongoDataSource() mongodbaccountdatasourceinterfaces.MemberAccessDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbaccountdatasourceinterfaces.MemberAccessDataSourceMongo) bool
}
