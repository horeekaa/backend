package databaseclientdependencies

import (
	container "github.com/golobby/container/v2"
	databasecoreclients "github.com/horeekaa/backend/core/databaseClient"
	databasecoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/interfaces/init"
	mongodbcoreclients "github.com/horeekaa/backend/core/databaseClient/mongoDB"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/transaction"
	mongodbcoreoperations "github.com/horeekaa/backend/core/databaseClient/mongoDB/operations"
	mongodbcoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/init"
	mongodbcoretransactions "github.com/horeekaa/backend/core/databaseClient/mongodb/transactions"
)

type DatabaseDependency struct{}

func (dbDependency *DatabaseDependency) bind() {
	container.Singleton(
		func() mongodbcoreclientinterfaces.MongoClient {
			mongoclient, _ := mongodbcoreclients.NewMongoClient()
			mongoclient.Connect()
			return mongoclient
		},
	)

	container.Singleton(
		func(mongoClient mongodbcoreclientinterfaces.MongoClient) databasecoreclientinterfaces.DatabaseClient {
			databaseClient, _ := databasecoreclients.NewDatabaseClient()
			databaseClient.SetMongoDBClient(mongoClient)
			return databaseClient
		},
	)

	container.Singleton(
		func(mongoClient mongodbcoreclientinterfaces.MongoClient) mongodbcoreoperationinterfaces.BasicOperation {
			basicOperation, _ := mongodbcoreoperations.NewBasicOperation(mongoClient)
			return basicOperation
		},
	)

	container.Transient(
		func(mongoClient mongodbcoreclientinterfaces.MongoClient) mongodbcoretransactioninterfaces.MongoRepoTransaction {
			mongoTransaction, _ := mongodbcoretransactions.NewMongoTransaction(mongoClient)
			return mongoTransaction
		},
	)
}