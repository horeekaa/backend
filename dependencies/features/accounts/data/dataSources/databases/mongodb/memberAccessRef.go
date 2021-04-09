package mongodbaccountdatasourcedependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	mongodbaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	databaseaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/sources"
)

type MemberAccessRefDataSourceDependency struct{}

func (memberAccessRefDataSourceDependency *MemberAccessRefDataSourceDependency) bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbaccountdatasourceinterfaces.MemberAccessRefDataSourceMongo {
			memberAccessRefRepoMongo, _ := mongodbaccountdatasources.NewMemberAccessRefDataSourceMongo(basicOperation)
			return memberAccessRefRepoMongo
		},
	)

	container.Singleton(
		func(memberAccessRefRepoMongo mongodbaccountdatasourceinterfaces.MemberAccessRefDataSourceMongo) databaseaccountdatasourceinterfaces.MemberAccessRefDataSource {
			memberAccessRefRepo, _ := databaseaccountdatasources.NewMemberAccessRefDataSource()
			memberAccessRefRepo.SetMongoDataSource(memberAccessRefRepoMongo)
			return memberAccessRefRepo
		},
	)
}
