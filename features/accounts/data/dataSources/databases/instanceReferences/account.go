package databaseaccountinstancereferences

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type AccountRepo struct {
	Instance *mongodbaccountdatasourceinterfaces.AccountRepoMongo
}
