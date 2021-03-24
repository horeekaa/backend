package databaseaccountrepos

import (
	databaseaccountrepointerfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/repos"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type memberAccessRepo struct {
	memberAccessRepoMongo mongodbaccountdatasourceinterfaces.MemberAccessRepoMongo
}

func (mbAccRepo *memberAccessRepo) SetMongoRepo(mongoRepo mongodbaccountdatasourceinterfaces.MemberAccessRepoMongo) bool {
	mbAccRepo.memberAccessRepoMongo = mongoRepo
	return true
}

func (mbAccRepo *memberAccessRepo) GetMongoRepo() mongodbaccountdatasourceinterfaces.MemberAccessRepoMongo {
	return mbAccRepo.memberAccessRepoMongo
}

func NewMemberAccessRepo() (databaseaccountrepointerfaces.MemberAccessRepo, error) {
	return &memberAccessRepo{}, nil
}
