package databaseaccountrepointerfaces

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type MemberAccessRepo interface {
	GetMongoRepo() mongodbaccountdatasourceinterfaces.MemberAccessRepoMongo
	SetMongoRepo(mongoRepo mongodbaccountdatasourceinterfaces.MemberAccessRepoMongo) bool
}
