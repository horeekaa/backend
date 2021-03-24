package mongodbaccountrepodependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	databaseaccountrepointerfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/repos"
	mongodbaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	databaseaccountrepos "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/repos"
)

type MemberAccessRepoDependency struct{}

func (memberAccessRepoDependency *MemberAccessRepoDependency) bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbaccountdatasourceinterfaces.MemberAccessRepoMongo {
			memberAccessRepoMongo, _ := mongodbaccountdatasources.NewMemberAccessRepoMongo(basicOperation)
			return memberAccessRepoMongo
		},
	)

	container.Singleton(
		func(memberAccessRepoMongo mongodbaccountdatasourceinterfaces.MemberAccessRepoMongo) databaseaccountrepointerfaces.MemberAccessRepo {
			memberAccessRepo, _ := databaseaccountrepos.NewMemberAccessRepo()
			memberAccessRepo.SetMongoRepo(memberAccessRepoMongo)
			return memberAccessRepo
		},
	)
}
