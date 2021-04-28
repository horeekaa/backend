package mongodbmemberaccessrefdatasourcedependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	mongodbmemberaccessrefdatasources "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/mongodb"
	mongodbmemberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/mongodb/interfaces"
	databasememberaccessrefdatasources "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/sources"
)

type MemberAccessRefDataSourceDependency struct{}

func (memberAccessRefDataSourceDependency *MemberAccessRefDataSourceDependency) bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbmemberaccessrefdatasourceinterfaces.MemberAccessRefDataSourceMongo {
			memberAccessRefRepoMongo, _ := mongodbmemberaccessrefdatasources.NewMemberAccessRefDataSourceMongo(basicOperation)
			return memberAccessRefRepoMongo
		},
	)

	container.Singleton(
		func(memberAccessRefRepoMongo mongodbmemberaccessrefdatasourceinterfaces.MemberAccessRefDataSourceMongo) databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource {
			memberAccessRefRepo, _ := databasememberaccessrefdatasources.NewMemberAccessRefDataSource()
			memberAccessRefRepo.SetMongoDataSource(memberAccessRefRepoMongo)
			return memberAccessRefRepo
		},
	)
}
