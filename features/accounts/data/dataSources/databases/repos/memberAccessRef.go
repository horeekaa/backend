package databaseaccountrepos

import (
	databaseaccountrepointerfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/repos"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type memberAccessRefRepo struct {
	memberAccessRefRepoMongo mongodbaccountdatasourceinterfaces.MemberAccessRefRepoMongo
}

func (mbAccRefRepo *memberAccessRefRepo) SetMongoRepo(mongoRepo mongodbaccountdatasourceinterfaces.MemberAccessRefRepoMongo) bool {
	mbAccRefRepo.memberAccessRefRepoMongo = mongoRepo
	return true
}

func (mbAccRefRepo *memberAccessRefRepo) GetMongoRepo() mongodbaccountdatasourceinterfaces.MemberAccessRefRepoMongo {
	return mbAccRefRepo.memberAccessRefRepoMongo
}

func NewMemberAccessRefRepo() (databaseaccountrepointerfaces.MemberAccessRefRepo, error) {
	return &memberAccessRefRepo{}, nil
}
