package databasememberaccessdatasourceinterfaces

import (
	mongodbmemberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/mongodb/interfaces"
)

type MemberAccessDataSource interface {
	GetMongoDataSource() mongodbmemberaccessdatasourceinterfaces.MemberAccessDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbmemberaccessdatasourceinterfaces.MemberAccessDataSourceMongo) bool
}
