package databaseaccountdatasources

import (
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type personDataSource struct {
	personDataSourceMongo mongodbaccountdatasourceinterfaces.PersonDataSourceMongo
}

func (personDataSource *personDataSource) SetMongoDataSource(mongoDataSource mongodbaccountdatasourceinterfaces.PersonDataSourceMongo) bool {
	personDataSource.personDataSourceMongo = mongoDataSource
	return true
}

func (personDataSource *personDataSource) GetMongoDataSource() mongodbaccountdatasourceinterfaces.PersonDataSourceMongo {
	return personDataSource.personDataSourceMongo
}

func NewPersonDataSource() (databaseaccountdatasourceinterfaces.PersonDataSource, error) {
	return &personDataSource{}, nil
}
