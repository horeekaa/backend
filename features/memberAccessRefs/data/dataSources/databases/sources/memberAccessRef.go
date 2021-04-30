package databasememberaccessrefdatasources

import (
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	mongodbmemberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/mongodb/interfaces"
)

type memberAccessRefDataSource struct {
	memberAccessRefDataSourceMongo mongodbmemberaccessrefdatasourceinterfaces.MemberAccessRefDataSourceMongo
}

func (mbAccRefDataSource *memberAccessRefDataSource) SetMongoDataSource(mongoDataSource mongodbmemberaccessrefdatasourceinterfaces.MemberAccessRefDataSourceMongo) bool {
	mbAccRefDataSource.memberAccessRefDataSourceMongo = mongoDataSource
	return true
}

func (mbAccRefDataSource *memberAccessRefDataSource) GetMongoDataSource() mongodbmemberaccessrefdatasourceinterfaces.MemberAccessRefDataSourceMongo {
	return mbAccRefDataSource.memberAccessRefDataSourceMongo
}

func NewMemberAccessRefDataSource() (databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource, error) {
	return &memberAccessRefDataSource{}, nil
}
