package databasecoreclients

import (
	databasecoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/interfaces/init"
	mongodbcoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/init"
)

type databaseClient struct {
	mongoClient mongodbcoreclientinterfaces.MongoClient
}

func (dbClient *databaseClient) SetMongoDBClient(mongoClient mongodbcoreclientinterfaces.MongoClient) bool {
	dbClient.mongoClient = mongoClient
	return true
}

func (dbClient *databaseClient) GetMongoDBClient() mongodbcoreclientinterfaces.MongoClient {
	return dbClient.mongoClient
}

func NewDatabaseClient() (databasecoreclientinterfaces.DatabaseClient, error) {
	return &databaseClient{}, nil
}
