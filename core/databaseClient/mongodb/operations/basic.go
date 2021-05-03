package mongodbcoreoperations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongodbcoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/init"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/repoExceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/repoExceptions/enums"
)

type basicOperation struct {
	mongoClient    mongodbcoreclientinterfaces.MongoClient
	collectionRef  *mongo.Collection
	collectionName string
	timeout        time.Duration
}

func NewBasicOperation(mongoClient mongodbcoreclientinterfaces.MongoClient) (mongodbcoreoperationinterfaces.BasicOperation, error) {
	return &basicOperation{
		mongoClient: mongoClient,
	}, nil
}

func (bscOperation *basicOperation) SetCollection(collectionName string) bool {
	client, _ := bscOperation.mongoClient.GetClient()
	databaseName, _ := bscOperation.mongoClient.GetDatabaseName()
	timeout, _ := bscOperation.mongoClient.GetDatabaseTimeout()

	bscOperation.collectionRef = client.Database(databaseName).Collection(collectionName)
	bscOperation.collectionName = collectionName
	bscOperation.timeout = timeout

	return true
}

func (bscOperation *basicOperation) GetCollectionName() string {
	return bscOperation.collectionName
}

func (bscOperation *basicOperation) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*time.Second)
	defer cancel()

	var res *mongo.SingleResult
	if &(*operationOptions).Session != nil {
		res = bscOperation.collectionRef.FindOne(*operationOptions.Session, bson.M{"_id": ID})
	} else {
		res = bscOperation.collectionRef.FindOne(ctx, bson.M{"_id": ID})
	}

	var objectReturn interface{}
	res.Decode(&objectReturn)
	if &objectReturn == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.IDNotFound,
			fmt.Sprintf("/%s/find", bscOperation.collectionName),
			errors.New(horeekaacoreexceptionenums.IDNotFound),
		)
	}

	return res, nil
}

func (bscOperation *basicOperation) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*time.Second)
	defer cancel()

	var bsonFilter bson.M
	data, err := bson.Marshal(query)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fmt.Sprintf("/%s/findOne", bscOperation.collectionName),
			err,
		)
	}
	bson.Unmarshal(data, &bsonFilter)

	var res *mongo.SingleResult
	if &(*operationOptions).Session != nil {
		res = bscOperation.collectionRef.FindOne(*operationOptions.Session, bsonFilter)
	} else {
		res = bscOperation.collectionRef.FindOne(ctx, bsonFilter)
	}

	return res, nil
}

func (bscOperation *basicOperation) Find(
	query map[string]interface{},
	paginationOpt *mongodbcoretypes.PaginationOptions,
	cursorDecoder func(cursorObject *mongo.Cursor) (interface{}, error),
	operationOptions *mongodbcoretypes.OperationOptions,
) (*bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*2*time.Second)
	defer cancel()

	var bsonFilter bson.M
	data, err := bson.Marshal(query)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fmt.Sprintf("/%s/find", bscOperation.collectionName),
			err,
		)
	}
	bson.Unmarshal(data, &bsonFilter)

	if paginationOpt.LastObjectID != nil {
		data, err = bson.Marshal(
			bson.M{
				"_id": bson.M{"$gt": *paginationOpt.LastObjectID},
			},
		)
		if err != nil {
			return nil, horeekaacoreexception.NewExceptionObject(
				horeekaacoreexceptionenums.UpstreamException,
				fmt.Sprintf("/%s/find", bscOperation.collectionName),
				err,
			)
		}
		bson.Unmarshal(data, &bsonFilter)
	}

	var opts *options.FindOptions
	opts.SetLimit(int64(10))
	if paginationOpt.QueryLimit != nil {
		opts.SetSort(bson.M{"_id": -1})
		opts.SetLimit(int64(*paginationOpt.QueryLimit))
	}

	var curr *mongo.Cursor
	if &(*operationOptions).Session != nil {
		curr, err = bscOperation.collectionRef.Find(*operationOptions.Session, bsonFilter, opts)
	} else {
		curr, err = bscOperation.collectionRef.Find(ctx, bsonFilter, opts)
	}

	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			fmt.Sprintf("/%s/find", bscOperation.collectionName),
			err,
		)
	}

	for curr.Next(ctx) {
		_, err := cursorDecoder(curr)
		if err != nil {
			return nil, horeekaacoreexception.NewExceptionObject(
				horeekaacoreexceptionenums.QueryObjectFailed,
				fmt.Sprintf("/%s/find", bscOperation.collectionName),
				err,
			)
		}
	}

	var output *bool
	*output = true
	return output, err
}

func (bscOperation *basicOperation) Create(input interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongodbcoretypes.CreateOperationOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*time.Second)
	defer cancel()

	var bsonObject bson.D
	data, err := bson.Marshal(input)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fmt.Sprintf("/%s/create", bscOperation.collectionName),
			err,
		)
	}
	bson.Unmarshal(data, &bsonObject)

	var res *mongo.InsertOneResult
	if &(*operationOptions).Session != nil {
		res, err = bscOperation.collectionRef.InsertOne(*operationOptions.Session, bsonObject)
	} else {
		res, err = bscOperation.collectionRef.InsertOne(ctx, bsonObject)
	}

	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.CreateObjectFailed,
			fmt.Sprintf("/%s/create", bscOperation.collectionName),
			err,
		)
	}

	return &mongodbcoretypes.CreateOperationOutput{
		ID:     res.InsertedID.(primitive.ObjectID),
		Object: input,
	}, nil
}

func (bscOperation *basicOperation) Update(ID primitive.ObjectID, updateData interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*time.Second)
	defer cancel()

	var bsonObject bson.D
	data, err := bson.Marshal(updateData)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fmt.Sprintf("/%s/update", bscOperation.collectionName),
			err,
		)
	}
	bson.Unmarshal(data, &bsonObject)

	if &(*operationOptions).Session != nil {
		_, err = bscOperation.collectionRef.UpdateOne(
			*operationOptions.Session,
			bson.M{"_id": ID},
			bson.D{primitive.E{Key: "$set", Value: bsonObject}},
		)
	} else {
		_, err = bscOperation.collectionRef.UpdateOne(
			ctx,
			bson.M{"_id": ID},
			bson.D{primitive.E{Key: "$set", Value: bsonObject}},
		)
	}

	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpdateObjectFailed,
			fmt.Sprintf("/%s/update", bscOperation.collectionName),
			err,
		)
	}
	res, err := bscOperation.FindByID(ID, operationOptions)

	return res, err
}
