package mongodbaccountrepodependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	databaseaccountrepointerfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/repos"
	mongodbaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	databaseaccountrepos "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/repos"
)

type PersonRepoDependency struct{}

func (personRepoDependency *PersonRepoDependency) bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbaccountdatasourceinterfaces.PersonRepoMongo {
			personRepoMongo, _ := mongodbaccountdatasources.NewPersonRepoMongo(basicOperation)
			return personRepoMongo
		},
	)

	container.Singleton(
		func(personRepoMongo mongodbaccountdatasourceinterfaces.PersonRepoMongo) databaseaccountrepointerfaces.PersonRepo {
			personRepo, _ := databaseaccountrepos.NewPersonRepo()
			personRepo.SetMongoRepo(personRepoMongo)
			return personRepo
		},
	)
}
