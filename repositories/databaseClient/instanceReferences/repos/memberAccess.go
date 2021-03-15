package databaseinstancereferences

import (
	mongorepointerfaces "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/interfaces/repos"
)

type MemberAccessRepo struct {
	Instance *mongorepointerfaces.MemberAccessRepoMongo
}
