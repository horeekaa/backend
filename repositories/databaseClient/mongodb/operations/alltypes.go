package mongooperations

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type OperationOptions struct {
	session *mongo.SessionContext
}
