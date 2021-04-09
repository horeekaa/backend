package mongodbaccountdatasourcedependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	mongodbaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	databaseaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/sources"
)

type PersonDataSourceDependency struct{}

func (personDataSourceDependency *PersonDataSourceDependency) bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbaccountdatasourceinterfaces.PersonDataSourceMongo {
			personRepoMongo, _ := mongodbaccountdatasources.NewPersonDataSourceMongo(basicOperation)
			return personRepoMongo
		},
	)

	container.Singleton(
		func(personRepoMongo mongodbaccountdatasourceinterfaces.PersonDataSourceMongo) databaseaccountdatasourceinterfaces.PersonDataSource {
			personRepo, _ := databaseaccountdatasources.NewPersonDataSource()
			personRepo.SetMongoDataSource(personRepoMongo)
			return personRepo
		},
	)
}
