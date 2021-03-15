package databaseinstancereferences

import (
	mongorepointerfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/repos"
)

type MemberAccessRefRepo struct {
	Instance *mongorepointerfaces.MemberAccessRefRepoMongo
}
