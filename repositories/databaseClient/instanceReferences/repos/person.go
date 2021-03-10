package databaseinstancereferences

import (
	mongorepointerfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/repos"
)

type PersonRepo struct {
	Instance *mongorepointerfaces.PersonRepoMongo
}
