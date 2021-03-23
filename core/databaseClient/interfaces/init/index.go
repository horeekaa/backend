package databasecoreclientinterfaces

import (
	mongodbcoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/init"
)

type DatabaseClient interface {
	SetMongoDBClient(mongoClient mongodbcoreclientinterfaces.MongoClient) bool
	GetMongoDBClient() mongodbcoreclientinterfaces.MongoClient
}
