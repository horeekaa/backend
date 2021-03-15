package databaseinstancereferences

import (
	mongorepointerfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/repos"
)

type AccountRepo struct {
	Instance *mongorepointerfaces.AccountRepoMongo
}
