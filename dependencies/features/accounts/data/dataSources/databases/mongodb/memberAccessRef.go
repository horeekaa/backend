package mongodbaccountrepodependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	databaseaccountrepointerfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/repos"
	mongodbaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	databaseaccountrepos "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/repos"
)

type MemberAccessRefRepoDependency struct{}

func (memberAccessRefRepoDependency *MemberAccessRefRepoDependency) bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbaccountdatasourceinterfaces.MemberAccessRefRepoMongo {
			memberAccessRefRepoMongo, _ := mongodbaccountdatasources.NewMemberAccessRefRepoMongo(basicOperation)
			return memberAccessRefRepoMongo
		},
	)

	container.Singleton(
		func(memberAccessRefRepoMongo mongodbaccountdatasourceinterfaces.MemberAccessRefRepoMongo) databaseaccountrepointerfaces.MemberAccessRefRepo {
			memberAccessRefRepo, _ := databaseaccountrepos.NewMemberAccessRefRepo()
			memberAccessRefRepo.SetMongoRepo(memberAccessRefRepoMongo)
			return memberAccessRefRepo
		},
	)
}
