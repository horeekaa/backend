package mongodbcoreoperations

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongodbcoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/init"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
)

type basicOperation struct {
	mongoClient         mongodbcoreclientinterfaces.MongoClient
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility
	collectionRef       mongodbcorewrapperinterfaces.MongoCollectionRef
	collectionName      string
	timeout             time.Duration
}

func NewBasicOperation(
	mongoClient mongodbcoreclientinterfaces.MongoClient,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (mongodbcoreoperationinterfaces.BasicOperation, error) {
	return &basicOperation{
		mongoClient:         mongoClient,
		mapProcessorUtility: mapProcessorUtility,
	}, nil
}

func (bscOperation *basicOperation) SetCollection(collectionName string) bool {
	timeout, _ := bscOperation.mongoClient.GetDatabaseTimeout()
	colRef, _ := bscOperation.mongoClient.GetCollectionRef(collectionName)

	bscOperation.collectionRef = colRef
	bscOperation.collectionName = collectionName
	bscOperation.timeout = timeout

	return true
}

func (bscOperation *basicOperation) GetCollectionName() string {
	return bscOperation.collectionName
}

func (bscOperation *basicOperation) FindByID(ID primitive.ObjectID, output interface{}, operationOptions *mongodbcoretypes.OperationOptions) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*time.Second)
	defer cancel()

	var res mongodbcorewrapperinterfaces.MongoSingleResult
	if operationOptions.Session != nil {
		res = bscOperation.collectionRef.FindOne(operationOptions.Session, bson.M{"_id": ID})
	} else {
		res = bscOperation.collectionRef.FindOne(ctx, bson.M{"_id": ID})
	}

	err := res.Decode(output)
	if err == mongo.ErrNoDocuments {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.IDNotFound,
			fmt.Sprintf("/%s/findByID", bscOperation.collectionName),
			nil,
		)
	}
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fmt.Sprintf("/%s/findByID", bscOperation.collectionName),
			nil,
		)
	}

	return true, nil
}

func (bscOperation *basicOperation) FindOne(query map[string]interface{}, output interface{}, operationOptions *mongodbcoretypes.OperationOptions) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*time.Second)
	defer cancel()

	bscOperation.mapProcessorUtility.RemoveNil(query)
	flattenedQuery := make(map[string]interface{})
	bscOperation.mapProcessorUtility.FlattenMap(
		"",
		query,
		&flattenedQuery,
	)

	var bsonFilter bson.M
	data, err := bson.Marshal(flattenedQuery)
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fmt.Sprintf("/%s/findOne", bscOperation.collectionName),
			err,
		)
	}
	bson.Unmarshal(data, &bsonFilter)

	var res mongodbcorewrapperinterfaces.MongoSingleResult
	if operationOptions.Session != nil {
		res = bscOperation.collectionRef.FindOne(operationOptions.Session, bsonFilter)
	} else {
		res = bscOperation.collectionRef.FindOne(ctx, bsonFilter)
	}

	err = res.Decode(output)
	if err == mongo.ErrNoDocuments {
		return true, err
	}
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fmt.Sprintf("/%s/findOne", bscOperation.collectionName),
			nil,
		)
	}

	return true, nil
}

func (bscOperation *basicOperation) Find(
	query map[string]interface{},
	paginationOpt *mongodbcoretypes.PaginationOptions,
	output interface{},
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*20*time.Second)
	defer cancel()

	bscOperation.mapProcessorUtility.RemoveNil(query)
	flattenedQuery := make(map[string]interface{})
	bscOperation.mapProcessorUtility.FlattenMap(
		"",
		query,
		&flattenedQuery,
	)

	var bsonFilter bson.M
	data, err := bson.Marshal(flattenedQuery)
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
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
			return false, horeekaacoreexception.NewExceptionObject(
				horeekaacoreexceptionenums.UpstreamException,
				fmt.Sprintf("/%s/find", bscOperation.collectionName),
				err,
			)
		}
		bson.Unmarshal(data, &bsonFilter)
	}

	opts := &options.FindOptions{}
	opts.SetLimit(int64(10))
	if paginationOpt.QueryLimit != nil {
		opts.SetSort(bson.M{"_id": 1})
		opts.SetLimit(int64(*paginationOpt.QueryLimit))
	}

	var curr mongodbcorewrapperinterfaces.MongoCursor
	if operationOptions.Session != nil {
		curr, err = bscOperation.collectionRef.Find(operationOptions.Session, bsonFilter, opts)
	} else {
		curr, err = bscOperation.collectionRef.Find(ctx, bsonFilter, opts)
	}

	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			fmt.Sprintf("/%s/find", bscOperation.collectionName),
			err,
		)
	}

	var outputList []*interface{}
	for curr.Next(ctx) {
		var outputSingle interface{}
		err := curr.Decode(&outputSingle)
		if err != nil {
			return false, horeekaacoreexception.NewExceptionObject(
				horeekaacoreexceptionenums.QueryObjectFailed,
				fmt.Sprintf("/%s/find", bscOperation.collectionName),
				err,
			)
		}

		outputList = append(outputList, &outputSingle)
	}
	output = outputList

	return true, err
}

func (bscOperation *basicOperation) Create(input interface{}, output interface{}, operationOptions *mongodbcoretypes.OperationOptions) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*time.Second)
	defer cancel()

	var createDataMap map[string]interface{}
	bsonTemp, err := bson.Marshal(input)
	bson.Unmarshal(bsonTemp, &createDataMap)
	bscOperation.mapProcessorUtility.RemoveNil(createDataMap)

	var bsonObject bson.M
	data, err := bson.Marshal(createDataMap)
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fmt.Sprintf("/%s/create", bscOperation.collectionName),
			err,
		)
	}
	bson.Unmarshal(data, &bsonObject)

	var res mongodbcorewrapperinterfaces.MongoInsertOneResult
	if operationOptions.Session != nil {
		res, err = bscOperation.collectionRef.InsertOne(operationOptions.Session, bsonObject)
	} else {
		res, err = bscOperation.collectionRef.InsertOne(ctx, bsonObject)
	}

	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.CreateObjectFailed,
			fmt.Sprintf("/%s/create", bscOperation.collectionName),
			err,
		)
	}

	insertedObjectID := map[string]interface{}{
		"_id": res.GetInsertedID().(primitive.ObjectID),
	}
	bsonTemp, _ = bson.Marshal(insertedObjectID)
	bson.Unmarshal(bsonTemp, output)

	bsonTemp, _ = bson.Marshal(input)
	bson.Unmarshal(bsonTemp, output)

	return true, nil
}

func (bscOperation *basicOperation) Update(ID primitive.ObjectID, updateData interface{}, output interface{}, operationOptions *mongodbcoretypes.OperationOptions) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bscOperation.timeout*time.Second)
	defer cancel()

	var updateDataMap map[string]interface{}
	bsonTemp, err := bson.Marshal(updateData)
	bson.Unmarshal(bsonTemp, &updateDataMap)
	bscOperation.mapProcessorUtility.RemoveNil(updateDataMap)

	var bsonObject bson.M
	data, err := bson.Marshal(updateDataMap)
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fmt.Sprintf("/%s/update", bscOperation.collectionName),
			err,
		)
	}
	bson.Unmarshal(data, &bsonObject)
	delete(bsonObject, "_id")

	if operationOptions.Session != nil {
		_, err = bscOperation.collectionRef.UpdateOne(
			operationOptions.Session,
			bson.M{"_id": ID},
			bson.M{"$set": bsonObject},
		)
	} else {
		_, err = bscOperation.collectionRef.UpdateOne(
			ctx,
			bson.M{"_id": ID},
			bson.M{"$set": bsonObject},
		)
	}
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpdateObjectFailed,
			fmt.Sprintf("/%s/update", bscOperation.collectionName),
			err,
		)
	}

	_, err = bscOperation.FindByID(ID, output, operationOptions)
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.IDNotFound,
			fmt.Sprintf("/%s/update", bscOperation.collectionName),
			err,
		)
	}

	return true, err
}
