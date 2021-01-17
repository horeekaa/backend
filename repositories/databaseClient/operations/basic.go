package mongooperations

import (
	"context"
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

func (bscOperation *BasicOperation) FindByID(id string, operationOptions *OperationOptions) (*interface{}, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.QueryObjectFailed,
			fmt.Sprintf("/%s/findById", bscOperation.CollectionName),
			err,
		)
	}
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.Timeout*time.Second)
	defer cancel()

	var res *mongo.SingleResult
	if &(*operationOptions).session != nil {
		res = bscOperation.CollectionRef.FindOne(*operationOptions.session, bson.M{"_id": objectID})
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

	var res *mongo.SingleResult
	if &(*operationOptions).session != nil {
		res = bscOperation.CollectionRef.FindOne(*operationOptions.session, query)
	} else {
		res = bscOperation.CollectionRef.FindOne(ctx, query)
	}

	var object *interface{}
	res.Decode(object)

	return object, nil
}

func (bscOperation *BasicOperation) Find(query map[string]interface{}, operationOptions *OperationOptions) ([]*interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.Timeout*2*time.Second)
	defer cancel()

	var curr *mongo.Cursor
	var err error
	if &(*operationOptions).session != nil {
		curr, err = bscOperation.CollectionRef.Find(*operationOptions.session, query)
	} else {
		curr, err = bscOperation.CollectionRef.Find(ctx, query)
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

func (bscOperation *BasicOperation) Create(input interface{}, operationOptions *OperationOptions) (*interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.Timeout*time.Second)
	defer cancel()

	var res *mongo.InsertOneResult
	var err error
	if &(*operationOptions).session != nil {
		res, err = bscOperation.CollectionRef.InsertOne(*operationOptions.session, input)
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

	return map[string]interface{}{
		"Return": &struct {
			ID     string
			Object interface{}
		}{res.InsertedID.(primitive.ObjectID).Hex(), input},
	}["Return"].(*interface{}), nil
}

func (bscOperation *BasicOperation) Update(id string, updateData interface{}, operationOptions *OperationOptions) (*interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.Timeout*time.Second)
	defer cancel()

	var err error
	if &(*operationOptions).session != nil {
		_, err = bscOperation.CollectionRef.UpdateOne(
			*operationOptions.session,
			bson.M{"_id": id},
			updateData,
		)
	} else {
		_, err = bscOperation.CollectionRef.UpdateOne(
			ctx,
			bson.M{"_id": id},
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
	objectOutput, errOutput := bscOperation.FindByID(id, operationOptions)

	return objectOutput, errOutput
}
