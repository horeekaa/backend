package mongodbcoreoperations

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	horeekaacoreexception "github.com/horeekaa/backend/core/_errors/repoExceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/_errors/repoExceptions/_enums"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	mongodbcoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/init"
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

func (bscOperation *basicOperation) FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error) {
	objectID := ID.(primitive.ObjectID)
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*time.Second)
	defer cancel()

	var res *mongo.SingleResult
	if &(*operationOptions).Session != nil {
		res = bscOperation.collectionRef.FindOne(*operationOptions.Session, bson.M{"_id": objectID})
	} else {
		res = bscOperation.collectionRef.FindOne(ctx, bson.M{"_id": objectID})
	}
	return res, nil
}

func (bscOperation *basicOperation) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*time.Second)
	defer cancel()

	var bsonObject bson.M
	encodedJSON, _ := json.Marshal(query)
	_ = bson.Unmarshal(encodedJSON, &bsonObject)

	var res *mongo.SingleResult
	if &(*operationOptions).Session != nil {
		res = bscOperation.collectionRef.FindOne(*operationOptions.Session, bsonObject)
	} else {
		res = bscOperation.collectionRef.FindOne(ctx, bsonObject)
	}

	return res, nil
}

func (bscOperation *basicOperation) Find(query map[string]interface{}, cursorDecoder func(cursorObject *mongodbcoretypes.CursorObject) (interface{}, error), operationOptions *mongodbcoretypes.OperationOptions) (*bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*2*time.Second)
	defer cancel()

	var bsonObject bson.M
	encodedJSON, _ := json.Marshal(query)
	_ = bson.Unmarshal(encodedJSON, &bsonObject)

	var curr *mongo.Cursor
	var err error
	if &(*operationOptions).Session != nil {
		curr, err = bscOperation.collectionRef.Find(*operationOptions.Session, bsonObject)
	} else {
		curr, err = bscOperation.collectionRef.Find(ctx, bsonObject)
	}

	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			fmt.Sprintf("/%s/find", bscOperation.collectionName),
			err,
		)
	}

	for curr.Next(ctx) {
		_, err := cursorDecoder(
			&mongodbcoretypes.CursorObject{
				MongoFindCursor: curr,
			},
		)
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

	var res *mongo.InsertOneResult
	var err error
	if &(*operationOptions).Session != nil {
		res, err = bscOperation.collectionRef.InsertOne(*operationOptions.Session, input)
	} else {
		res, err = bscOperation.collectionRef.InsertOne(ctx, input)
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

func (bscOperation *basicOperation) Update(ID interface{}, updateData interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error) {
	objectID := ID.(primitive.ObjectID)
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*time.Second)
	defer cancel()

	var err error
	if &(*operationOptions).Session != nil {
		_, err = bscOperation.collectionRef.UpdateOne(
			*operationOptions.Session,
			bson.M{"_id": objectID},
			updateData,
		)
	} else {
		_, err = bscOperation.collectionRef.UpdateOne(
			ctx,
			bson.M{"_id": objectID},
			updateData,
		)
	}

	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpdateObjectFailed,
			fmt.Sprintf("/%s/update", bscOperation.collectionName),
			err,
		)
	}
	res, err := bscOperation.FindByID(objectID, operationOptions)

	return res, err
}
