package databasememberaccessdatasources

import (
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	mongodbmemberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/mongodb/interfaces"
)

type memberAccessDataSource struct {
	memberAccessDataSourceMongo mongodbmemberaccessdatasourceinterfaces.MemberAccessDataSourceMongo
}

func (mbAccDataSource *memberAccessDataSource) SetMongoDataSource(mongoDataSource mongodbmemberaccessdatasourceinterfaces.MemberAccessDataSourceMongo) bool {
	mbAccDataSource.memberAccessDataSourceMongo = mongoDataSource
	return true
}

func (mbAccDataSource *memberAccessDataSource) GetMongoDataSource() mongodbmemberaccessdatasourceinterfaces.MemberAccessDataSourceMongo {
	return mbAccDataSource.memberAccessDataSourceMongo
}

func NewMemberAccessDataSource() (databasememberaccessdatasourceinterfaces.MemberAccessDataSource, error) {
	return &memberAccessDataSource{}, nil
}
