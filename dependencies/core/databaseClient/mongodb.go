package databaseclientdependencies

import (
	container "github.com/golobby/container/v2"
	databaseclientinterfaces "github.com/horeekaa/backend/core/databaseClient/interfaces/init"
	mongodbcoreclients "github.com/horeekaa/backend/core/databaseClient/mongoDB"
)

type DatabaseDependency struct{}

func (dbDependency *DatabaseDependency) bind() {
	container.Singleton(
		func() databaseclientinterfaces.DatabaseClient {
			mongoclient, _ := mongodbcoreclients.NewMongoClient()
			return mongoclient
		},
	)

	container.Make(
		func(dbClient databaseclientinterfaces.DatabaseClient) {
			dbClient.Connect()
		},
	)
}
