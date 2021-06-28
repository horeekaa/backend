package mongodbcorewrapperinterfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSingleResult interface {
	Decode(v interface{}) error
}

type MongoCursor interface {
	Next(ctx context.Context) bool
	Decode(val interface{}) error
}

type MongoInsertOneResult interface {
	GetInsertedID() interface{}
}

type MongoCollectionRef interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) MongoSingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (MongoCursor, error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (MongoInsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (interface{}, error)
}
