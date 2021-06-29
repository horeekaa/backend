package mongodbmemberaccessdatasourcedependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	mongodbmemberaccessdatasources "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/mongodb"
	mongodbmemberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/mongodb/interfaces"
	databasememberaccessdatasources "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/sources"
)

type MemberAccessDataSourceDependency struct{}

func (memberAccessDataSourceDependency *MemberAccessDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbmemberaccessdatasourceinterfaces.MemberAccessDataSourceMongo {
			memberAccessRepoMongo, _ := mongodbmemberaccessdatasources.NewMemberAccessDataSourceMongo(basicOperation)
			return memberAccessRepoMongo
		},
	)

	container.Singleton(
		func(memberAccessRepoMongo mongodbmemberaccessdatasourceinterfaces.MemberAccessDataSourceMongo) databasememberaccessdatasourceinterfaces.MemberAccessDataSource {
			memberAccessRepo, _ := databasememberaccessdatasources.NewMemberAccessDataSource()
			memberAccessRepo.SetMongoDataSource(memberAccessRepoMongo)
			return memberAccessRepo
		},
	)
}
