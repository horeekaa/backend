package databaseaccountrepos

import (
	databaseaccountrepointerfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/repos"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type personRepo struct {
	personRepoMongo mongodbaccountdatasourceinterfaces.PersonRepoMongo
}

func (personRepo *personRepo) SetMongoRepo(mongoRepo mongodbaccountdatasourceinterfaces.PersonRepoMongo) bool {
	personRepo.personRepoMongo = mongoRepo
	return true
}

func (personRepo *personRepo) GetMongoRepo() mongodbaccountdatasourceinterfaces.PersonRepoMongo {
	return personRepo.personRepoMongo
}

func NewPersonRepo() (databaseaccountrepointerfaces.PersonRepo, error) {
	return &personRepo{}, nil
}
