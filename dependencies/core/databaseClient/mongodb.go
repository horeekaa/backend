package databaseclientdependencies

import (
	container "github.com/golobby/container/v2"
	databasecoreclients "github.com/horeekaa/backend/core/databaseClient"
	databasecoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/interfaces/init"
	mongodbcoreclients "github.com/horeekaa/backend/core/databaseClient/mongodb"
	mongodbcoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/init"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoreoperations "github.com/horeekaa/backend/core/databaseClient/mongodb/operations"
	mongodbcorequerybuilders "github.com/horeekaa/backend/core/databaseClient/mongodb/queryBuilders"
	mongodbcoretransactions "github.com/horeekaa/backend/core/databaseClient/mongodb/transactions"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
)

type DatabaseDependency struct{}

func (dbDependency *DatabaseDependency) Bind() {
	container.Singleton(
		func() mongodbcorequerybuilderinterfaces.MongoQueryBuilder {
			queryBuilder, _ := mongodbcorequerybuilders.NewMongoQueryBuilder()
			return queryBuilder
		},
	)

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

	container.Transient(
		func(
			mongoClient mongodbcoreclientinterfaces.MongoClient,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) mongodbcoreoperationinterfaces.BasicOperation {
			basicOperation, _ := mongodbcoreoperations.NewBasicOperation(
				mongoClient,
				mapProcessorUtility,
			)
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
