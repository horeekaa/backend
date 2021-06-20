package mongodbloggingdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
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

func (orgDataSourceMongo *loggingDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error) {
	var output model.Logging
	_, err := orgDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (orgDataSourceMongo *loggingDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error) {
	var output model.Logging
	_, err := orgDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (orgDataSourceMongo *loggingDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Logging, error) {
	var loggings = []*model.Logging{}
	appendingFn := func(dbItem interface{}) (bool, error) {
		logging := dbItem.(model.Logging)
		loggings = append(loggings, &logging)
		return true, nil
	}
	_, err := orgDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return loggings, err
}

func (orgDataSourceMongo *loggingDataSourceMongo) Create(input *model.CreateLogging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error) {
	defaultedInput, err := orgDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Logging
	_, err = orgDataSourceMongo.basicOperation.Create(*defaultedInput.CreateLogging, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (orgDataSourceMongo *loggingDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.UpdateLogging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error) {
	updateData.ID = ID
	defaultedInput, err := orgDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Logging
	_, err = orgDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateLogging, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type setLoggingDefaultValuesOutput struct {
	CreateLogging *model.CreateLogging
	UpdateLogging *model.UpdateLogging
}

func (orgDataSourceMongo *loggingDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setLoggingDefaultValuesOutput, error) {
	currentTime := time.Now()

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.UpdateLogging)
		_, err := orgDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setLoggingDefaultValuesOutput{
			UpdateLogging: &updateInput,
		}, nil
	}
	createInput := (input).(model.CreateLogging)
	createInput.CreatedAt = &currentTime

	return &setLoggingDefaultValuesOutput{
		CreateLogging: &createInput,
	}, nil
}
