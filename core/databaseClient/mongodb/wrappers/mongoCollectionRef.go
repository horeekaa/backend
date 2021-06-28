package mongodbcorewrappers

import (
	"context"

	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoCollectionRef struct {
	*mongo.Collection
}

func (colRef *mongoCollectionRef) InsertOne(
	ctx context.Context,
	document interface{},
	opts ...*options.InsertOneOptions,
) (mongodbcorewrapperinterfaces.MongoInsertOneResult, error) {
	res, err := colRef.Collection.InsertOne(ctx, document, opts...)
	if err != nil {
		return nil, err
	}

	return NewMongoInsertOneResult(res)
}

func (colRef *mongoCollectionRef) Find(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (mongodbcorewrapperinterfaces.MongoCursor, error) {
	cursor, err := colRef.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func (colRef *mongoCollectionRef) FindOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) mongodbcorewrapperinterfaces.MongoSingleResult {
	res := colRef.Collection.FindOne(ctx, filter, opts...)

	return res
}

func (colRef *mongoCollectionRef) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (interface{}, error) {
	out, err := colRef.Collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func NewMongoCollectionRef(wrappedMongoCollection *mongo.Collection) (mongodbcorewrapperinterfaces.MongoCollectionRef, error) {
	return &mongoCollectionRef{
		wrappedMongoCollection,
	}, nil
}
