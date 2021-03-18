package databaseaccountinstancereferences

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type MemberAccessRepo struct {
	Instance *mongodbaccountdatasourceinterfaces.MemberAccessRepoMongo
}
