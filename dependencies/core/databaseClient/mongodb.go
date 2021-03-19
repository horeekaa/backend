package databaseclientdependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreclients "github.com/horeekaa/backend/core/databaseClient/mongoDB"
	databaseclientinterfaces "github.com/horeekaa/backend/core/databaseClient/interfaces/init"
)

type DatabaseDependency struct{}

func (dbDependency *DatabaseDependency) bind() {
	container.Singleton(
		func() databaseclientinterfaces.DatabaseClient {
			mongodbclient, _ := mongodbcoreclients.NewMongoClient()
			return mongodbclient
		},
	)
}
