package mongodbaccountdatasourcedependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	mongodbaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	databaseaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/sources"
)

type MemberAccessDataSourceDependency struct{}

func (memberAccessDataSourceDependency *MemberAccessDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbaccountdatasourceinterfaces.MemberAccessDataSourceMongo {
			memberAccessRepoMongo, _ := mongodbaccountdatasources.NewMemberAccessDataSourceMongo(basicOperation)
			return memberAccessRepoMongo
		},
	)

	container.Singleton(
		func(memberAccessRepoMongo mongodbaccountdatasourceinterfaces.MemberAccessDataSourceMongo) databaseaccountdatasourceinterfaces.MemberAccessDataSource {
			memberAccessRepo, _ := databaseaccountdatasources.NewMemberAccessDataSource()
			memberAccessRepo.SetMongoDataSource(memberAccessRepoMongo)
			return memberAccessRepo
		},
	)
}
