package databaseaccountdatasources

import (
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type accountDataSource struct {
	accountDataSourceMongo mongodbaccountdatasourceinterfaces.AccountDataSourceMongo
}

func (accDataSource *accountDataSource) SetMongoDataSource(mongoDataSource mongodbaccountdatasourceinterfaces.AccountDataSourceMongo) bool {
	accDataSource.accountDataSourceMongo = mongoDataSource
	return true
}

func (accDataSource *accountDataSource) GetMongoDataSource() mongodbaccountdatasourceinterfaces.AccountDataSourceMongo {
	return accDataSource.accountDataSourceMongo
}

func NewAccountDataSource() (databaseaccountdatasourceinterfaces.AccountDataSource, error) {
	return &accountDataSource{}, nil
}
