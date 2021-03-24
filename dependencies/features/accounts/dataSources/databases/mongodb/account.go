package mongodbaccountrepodependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	databaseaccountrepointerfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/repos"
	mongodbaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	databaseaccountrepos "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/repos"
)

type AccountRepoDependency struct{}

func (accountRepoDependency *AccountRepoDependency) bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbaccountdatasourceinterfaces.AccountRepoMongo {
			accountRepoMongo, _ := mongodbaccountdatasources.NewAccountRepoMongo(basicOperation)
			return accountRepoMongo
		},
	)

	container.Singleton(
		func(accountRepoMongo mongodbaccountdatasourceinterfaces.AccountRepoMongo) databaseaccountrepointerfaces.AccountRepo {
			accountRepo, _ := databaseaccountrepos.NewAccountRepo()
			accountRepo.SetMongoRepo(accountRepoMongo)
			return accountRepo
		},
	)
}
