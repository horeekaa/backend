package mongodbloggingdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongoDB/interfaces/operations"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	mongodbloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
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

func (orgDataSourceMongo *loggingDataSourceMongo) FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error) {
	res, err := orgDataSourceMongo.basicOperation.FindByID(ID, operationOptions)
	var output model.Logging
	res.Decode(&output)
	return &output, err
}

func (orgDataSourceMongo *loggingDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error) {
	res, err := orgDataSourceMongo.basicOperation.FindOne(query, operationOptions)
	var output model.Logging
	res.Decode(&output)
	return &output, err
}

func (orgDataSourceMongo *loggingDataSourceMongo) Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.Logging, error) {
	var loggings = []*model.Logging{}
	cursorDecoder := func(cursor *mongodbcoretypes.CursorObject) (interface{}, error) {
		var logging *model.Logging
		err := cursor.MongoFindCursor.Decode(logging)
		if err != nil {
			return nil, err
		}
		loggings = append(loggings, logging)
		return nil, nil
	}

	_, err := orgDataSourceMongo.basicOperation.Find(query, cursorDecoder, operationOptions)
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

	output, err := orgDataSourceMongo.basicOperation.Create(*defaultedInput.CreateLogging, operationOptions)
	if err != nil {
		return nil, err
	}

	loggingOutput := output.Object.(model.Logging)
	loggingOutput.ID = output.ID

	return &loggingOutput, err
}

func (orgDataSourceMongo *loggingDataSourceMongo) Update(ID interface{}, updateData *model.UpdateLogging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error) {
	defaultedInput, err := orgDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	res, err := orgDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateLogging, operationOptions)
	var output model.Logging
	res.Decode(&output)

	return &output, err
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