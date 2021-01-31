package mongooperations

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	horeekaaexception "github.com/horeekaa/backend/_errors/repoExceptions"
	horeekaaexceptionenums "github.com/horeekaa/backend/_errors/repoExceptions/_enums"
)

type BasicOperation struct {
	Client         *mongo.Client
	CollectionRef  *mongo.Collection
	Timeout        time.Duration
	CollectionName string
}

func (bscOperation *BasicOperation) FindByID(ID interface{}, operationOptions *OperationOptions) (*interface{}, error) {
	objectID := ID.(primitive.ObjectID)
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.Timeout*time.Second)
	defer cancel()

	var res *mongo.SingleResult
	if &(*operationOptions).Session != nil {
		res = bscOperation.CollectionRef.FindOne(*operationOptions.Session, bson.M{"_id": objectID})
	} else {
		res = bscOperation.CollectionRef.FindOne(ctx, bson.M{"_id": objectID})
	}

	var object *interface{}
	res.Decode(object)

	return object, nil
}

func (bscOperation *BasicOperation) FindOne(query map[string]interface{}, operationOptions *OperationOptions) (*interface{}, error) {
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

	var object *interface{}
	res.Decode(object)

	return object, nil
}

func (bscOperation *BasicOperation) Find(query map[string]interface{}, operationOptions *OperationOptions) ([]*interface{}, error) {
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
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.QueryObjectFailed,
			fmt.Sprintf("/%s/find", bscOperation.CollectionName),
			err,
		)
	}

	var objects []*interface{}
	for curr.Next(ctx) {
		var object *interface{}
		err := curr.Decode(object)
		if err != nil {
			return nil, horeekaaexception.NewExceptionObject(
				horeekaaexceptionenums.QueryObjectFailed,
				fmt.Sprintf("/%s/find", bscOperation.CollectionName),
				err,
			)
		}
		objects = append(objects, object)
	}

	return objects, err
}

func (bscOperation *BasicOperation) Create(input interface{}, operationOptions *OperationOptions) (*CreateOperationOutput, error) {
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
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.CreateObjectFailed,
			fmt.Sprintf("/%s/create", bscOperation.CollectionName),
			err,
		)
	}

	return &CreateOperationOutput{
		ID:     res.InsertedID.(primitive.ObjectID),
		Object: input,
	}, nil
}

func (bscOperation *BasicOperation) Update(ID interface{}, updateData interface{}, operationOptions *OperationOptions) (*interface{}, error) {
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
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.UpdateObjectFailed,
			fmt.Sprintf("/%s/update", bscOperation.CollectionName),
			err,
		)
	}
	objectOutput, errOutput := bscOperation.FindByID(objectID, operationOptions)

	return objectOutput, errOutput
}
