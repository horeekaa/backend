package databaseaccountrepointerfaces

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type AccountRepo interface {
	GetMongoRepo() mongodbaccountdatasourceinterfaces.AccountRepoMongo
	SetMongoRepo(mongoRepo mongodbaccountdatasourceinterfaces.AccountRepoMongo) bool
}
