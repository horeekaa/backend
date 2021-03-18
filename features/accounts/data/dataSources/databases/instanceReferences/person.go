package databaseaccountinstancereferences

import (
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type PersonRepo struct {
	Instance *mongodbaccountdatasourceinterfaces.PersonRepoMongo
}
