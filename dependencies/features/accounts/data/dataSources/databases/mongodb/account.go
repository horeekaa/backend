package mongodbaccountdatasourcedependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	mongodbaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb"
	mongodbaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/mongodb/interfaces"
	databaseaccountdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/sources"
)

type AccountDataSourceDependency struct{}

func (accountDataSourceDependency *AccountDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbaccountdatasourceinterfaces.AccountDataSourceMongo {
			accountDataSourceMongo, _ := mongodbaccountdatasources.NewAccountDataSourceMongo(basicOperation)
			return accountDataSourceMongo
		},
	)

	container.Singleton(
		func(accountDataSourceMongo mongodbaccountdatasourceinterfaces.AccountDataSourceMongo) databaseaccountdatasourceinterfaces.AccountDataSource {
			accountRepo, _ := databaseaccountdatasources.NewAccountDataSource()
			accountRepo.SetMongoDataSource(accountDataSourceMongo)
			return accountRepo
		},
	)
}
