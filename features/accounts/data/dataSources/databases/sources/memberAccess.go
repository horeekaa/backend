package databaseaccountdatasources

import (
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
)

type memberAccessDataSource struct {
	memberAccessDataSourceMongo mongodbaccountdatasourceinterfaces.MemberAccessDataSourceMongo
}

func (mbAccDataSource *memberAccessDataSource) SetMongoDataSource(mongoDataSource mongodbaccountdatasourceinterfaces.MemberAccessDataSourceMongo) bool {
	mbAccDataSource.memberAccessDataSourceMongo = mongoDataSource
	return true
}

func (mbAccDataSource *memberAccessDataSource) GetMongoDataSource() mongodbaccountdatasourceinterfaces.MemberAccessDataSourceMongo {
	return mbAccDataSource.memberAccessDataSourceMongo
}

func NewMemberAccessDataSource() (databaseaccountdatasourceinterfaces.MemberAccessDataSource, error) {
	return &memberAccessDataSource{}, nil
}
