package databaseaccountinstancereferences

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type MemberAccessRefRepo struct {
	Instance *mongodbaccountdatasourceinterfaces.MemberAccessRefRepoMongo
}
