package mongodbloggingdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoggingDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Logging, error)
	Create(input *model.CreateLogging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error)
	Update(ID primitive.ObjectID, updateData *model.UpdateLogging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Logging, error)
}
