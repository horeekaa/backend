package mongodbloggingdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type loggingDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewLoggingDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbloggingdatasourceinterfaces.LoggingDataSourceMongo, error) {
	basicOperation.SetCollection("loggings")
	return &loggingDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (logDataSourceMongo *loggingDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error) {
	var output model.Logging
	_, err := logDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (logDataSourceMongo *loggingDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error) {
	var output model.Logging
	_, err := logDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (logDataSourceMongo *loggingDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Logging, error) {
	var loggings = []*model.Logging{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var logging model.Logging
		if err := cursor.Decode(&logging); err != nil {
			return err
		}
		loggings = append(loggings, &logging)
		return nil
	}
	_, err := logDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return loggings, err
}

func (logDataSourceMongo *loggingDataSourceMongo) Create(input *model.CreateLogging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error) {
	_, err := logDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Logging
	_, err = logDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (logDataSourceMongo *loggingDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.UpdateLogging,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.Logging, error) {
	_, err := logDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Logging
	_, err = logDataSourceMongo.basicOperation.Update(
		updateCriteria,
		map[string]interface{}{
			"$set": updateData,
		},
		&output,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (logDataSourceMongo *loggingDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.UpdateLogging,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()
	existingObject, err := logDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			"/loggingDataSource/update",
			nil,
		)
	}
	input.UpdatedAt = &currentTime

	return true, nil
}

func (logDataSourceMongo *loggingDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.CreateLogging,
) (bool, error) {
	currentTime := time.Now()
	input.CreatedAt = &currentTime

	return true, nil
}
