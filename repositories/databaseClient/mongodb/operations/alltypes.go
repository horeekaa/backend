package mongooperations

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OperationOptions struct {
	Session *mongo.SessionContext
}

type CreateOperationOutput struct {
	ID     primitive.ObjectID
	Object interface{}
}
