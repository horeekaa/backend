package mongooperations

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OperationOptions struct {
	session *mongo.SessionContext
}

type OperationQueryType bson.M

type CreateOperationOutput struct {
	ID     primitive.ObjectID
	Object interface{}
}
