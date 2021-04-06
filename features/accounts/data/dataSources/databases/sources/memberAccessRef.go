package databaseaccountdatasources

import (
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type memberAccessRefDataSource struct {
	memberAccessRefDataSourceMongo mongodbaccountdatasourceinterfaces.MemberAccessRefDataSourceMongo
}

func (mbAccRefDataSource *memberAccessRefDataSource) SetMongoDataSource(mongoDataSource mongodbaccountdatasourceinterfaces.MemberAccessRefDataSourceMongo) bool {
	mbAccRefDataSource.memberAccessRefDataSourceMongo = mongoDataSource
	return true
}

func (mbAccRefDataSource *memberAccessRefDataSource) GetMongoDataSource() mongodbaccountdatasourceinterfaces.MemberAccessRefDataSourceMongo {
	return mbAccRefDataSource.memberAccessRefDataSourceMongo
}

func NewMemberAccessRefDataSource() (databaseaccountdatasourceinterfaces.MemberAccessRefDataSource, error) {
	return &memberAccessRefDataSource{}, nil
}
