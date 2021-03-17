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
	mongodbcoremodels "github.com/horeekaa/backend/core/databaseClient/mongoDB/operations/models"
)

type basicOperation struct {
	Client         *mongo.Client
	CollectionRef  *mongo.Collection
	Timeout        time.Duration
	CollectionName string
}

func NewBasicOperation(client *mongo.Client, collectionRef *mongo.Collection, timeout time.Duration, collectionName string) (mongodbcoreoperationinterfaces.BasicOperation, error) {
	return &basicOperation{
		Client:         client,
		CollectionRef:  collectionRef,
		Timeout:        timeout,
		CollectionName: collectionName,
	}, nil
}

func (bscOperation *basicOperation) FindByID(ID interface{}, operationOptions *mongodbcoremodels.OperationOptions) (*mongo.SingleResult, error) {
	objectID := ID.(primitive.ObjectID)
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.Timeout*time.Second)
	defer cancel()

	var res *mongo.SingleResult
	if &(*operationOptions).Session != nil {
		res = bscOperation.CollectionRef.FindOne(*operationOptions.Session, bson.M{"_id": objectID})
	} else {
		res = bscOperation.CollectionRef.FindOne(ctx, bson.M{"_id": objectID})
	}
	return res, nil
}

func (bscOperation *basicOperation) FindOne(query map[string]interface{}, operationOptions *mongodbcoremodels.OperationOptions) (*mongo.SingleResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.Timeout*time.Second)
	defer cancel()

	var bsonObject bson.M
	encodedJSON, _ := json.Marshal(query)
	_ = bson.Unmarshal(encodedJSON, &bsonObject)

	var res *mongo.SingleResult
	if &(*operationOptions).Session != nil {
		res = bscOperation.CollectionRef.FindOne(*operationOptions.Session, bsonObject)
	} else {
		res = bscOperation.CollectionRef.FindOne(ctx, bsonObject)
	}

	return res, nil
}

func (bscOperation *basicOperation) Find(query map[string]interface{}, cursorDecoder func(cursorObject *mongodbcoremodels.CursorObject) (interface{}, error), operationOptions *mongodbcoremodels.OperationOptions) (*bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.Timeout*2*time.Second)
	defer cancel()

	var bsonObject bson.M
	encodedJSON, _ := json.Marshal(query)
	_ = bson.Unmarshal(encodedJSON, &bsonObject)

	var curr *mongo.Cursor
	var err error
	if &(*operationOptions).Session != nil {
		curr, err = bscOperation.CollectionRef.Find(*operationOptions.Session, bsonObject)
	} else {
		curr, err = bscOperation.CollectionRef.Find(ctx, bsonObject)
	}

	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			fmt.Sprintf("/%s/find", bscOperation.CollectionName),
			err,
		)
	}

	for curr.Next(ctx) {
		_, err := cursorDecoder(
			&mongodbcoremodels.CursorObject{
				MongoFindCursor: curr,
			},
		)
		if err != nil {
			return nil, horeekaacoreexception.NewExceptionObject(
				horeekaacoreexceptionenums.QueryObjectFailed,
				fmt.Sprintf("/%s/find", bscOperation.CollectionName),
				err,
			)
		}
	}

	var output *bool
	*output = true
	return output, err
}

func (bscOperation *basicOperation) Create(input interface{}, operationOptions *mongodbcoremodels.OperationOptions) (*mongodbcoremodels.CreateOperationOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.Timeout*time.Second)
	defer cancel()

	var res *mongo.InsertOneResult
	var err error
	if &(*operationOptions).Session != nil {
		res, err = bscOperation.CollectionRef.InsertOne(*operationOptions.Session, input)
	} else {
		res, err = bscOperation.CollectionRef.InsertOne(ctx, input)
	}

	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.CreateObjectFailed,
			fmt.Sprintf("/%s/create", bscOperation.CollectionName),
			err,
		)
	}

	return &mongodbcoremodels.CreateOperationOutput{
		ID:     res.InsertedID.(primitive.ObjectID),
		Object: input,
	}, nil
}

func (bscOperation *basicOperation) Update(ID interface{}, updateData interface{}, operationOptions *mongodbcoremodels.OperationOptions) (*mongo.SingleResult, error) {
	objectID := ID.(primitive.ObjectID)
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.Timeout*time.Second)
	defer cancel()

	var err error
	if &(*operationOptions).Session != nil {
		_, err = bscOperation.CollectionRef.UpdateOne(
			*operationOptions.Session,
			bson.M{"_id": objectID},
			updateData,
		)
	} else {
		_, err = bscOperation.CollectionRef.UpdateOne(
			ctx,
			bson.M{"_id": objectID},
			updateData,
		)
	}

	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpdateObjectFailed,
			fmt.Sprintf("/%s/update", bscOperation.CollectionName),
			err,
		)
	}
	res, err := bscOperation.FindByID(objectID, operationOptions)

	return res, err
}
