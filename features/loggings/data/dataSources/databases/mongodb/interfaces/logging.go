package mongodbloggingdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	"github.com/horeekaa/backend/model"
)

type LoggingDataSourceMongo interface {
	FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error)
	Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.Logging, error)
	Create(input *model.CreateLogging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error)
	Update(ID interface{}, updateData *model.UpdateLogging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error)
}
