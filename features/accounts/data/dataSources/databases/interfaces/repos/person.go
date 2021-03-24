package databaseaccountrepointerfaces

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type PersonRepo interface {
	GetMongoRepo() mongodbaccountdatasourceinterfaces.PersonRepoMongo
	SetMongoRepo(mongoRepo mongodbaccountdatasourceinterfaces.PersonRepoMongo) bool
}
