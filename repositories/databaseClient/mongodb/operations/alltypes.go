package mongooperations

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OperationOptions struct {
	session *mongo.SessionContext
}

type OperationQueryType bson.M
