package mongodbcoretypes

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DefaultValuesUpdateType string = "UPDATE"
	DefaultValuesCreateType string = "CREATE"
)

type DefaultValuesOptions struct {
	DefaultValuesType string
}

type OperationOptions struct {
	Session *mongo.SessionContext
}

type CursorObject struct {
	MongoFindCursor *mongo.Cursor
}

type CreateOperationOutput struct {
	ID     primitive.ObjectID
	Object interface{}
}
