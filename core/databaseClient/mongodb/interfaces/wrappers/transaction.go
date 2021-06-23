package mongodbcorewrapperinterfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSession interface {
	AbortTransaction(_a0 context.Context) error
	AdvanceClusterTime(_a0 bson.Raw) error
	AdvanceOperationTime(_a0 *primitive.Timestamp) error
	Client() *mongo.Client
	ClusterTime() bson.Raw
	CommitTransaction(_a0 context.Context) error
	EndSession(_a0 context.Context)
	ID() bson.Raw
	OperationTime() *primitive.Timestamp
	StartTransaction(_a0 ...*options.TransactionOptions) error
	WithTransaction(ctx context.Context, fn func(mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error)
}

type MongoSessionContext interface {
	mongo.SessionContext
}
