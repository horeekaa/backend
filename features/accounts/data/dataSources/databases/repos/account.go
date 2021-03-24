package databaseaccountrepos

import (
	databaseaccountrepointerfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/repos"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type accountRepo struct {
	accountRepoMongo mongodbaccountdatasourceinterfaces.AccountRepoMongo
}

func (accRepo *accountRepo) SetMongoRepo(mongoRepo mongodbaccountdatasourceinterfaces.AccountRepoMongo) bool {
	accRepo.accountRepoMongo = mongoRepo
	return true
}

func (accRepo *accountRepo) GetMongoRepo() mongodbaccountdatasourceinterfaces.AccountRepoMongo {
	return accRepo.accountRepoMongo
}

func NewAccountRepo() (databaseaccountrepointerfaces.AccountRepo, error) {
	return &accountRepo{}, nil
}
