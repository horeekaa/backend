package databaseaccountrepointerfaces

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type MemberAccessRefRepo interface {
	GetMongoRepo() mongodbaccountdatasourceinterfaces.MemberAccessRefRepoMongo
	SetMongoRepo(mongoRepo mongodbaccountdatasourceinterfaces.MemberAccessRefRepoMongo) bool
}
