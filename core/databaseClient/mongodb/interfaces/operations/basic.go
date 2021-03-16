package mongodbcoreoperationinterfaces

import (
	mongodbcoremodels "github.com/horeekaa/backend/core/databaseClient/mongoDB/operations/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type BasicOperation interface {
	FindByID(ID interface{}, operationOptions *mongodbcoremodels.OperationOptions) (*mongo.SingleResult, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoremodels.OperationOptions) (*mongo.SingleResult, error)
	Find(query map[string]interface{}, cursorDecoder func(cursorObject *mongodbcoremodels.CursorObject) (interface{}, error), operationOptions *mongodbcoremodels.OperationOptions) (*bool, error)
	Create(input interface{}, operationOptions *mongodbcoremodels.OperationOptions) (*mongodbcoremodels.CreateOperationOutput, error)
	Update(ID interface{}, updateData interface{}, operationOptions *mongodbcoremodels.OperationOptions) (*mongo.SingleResult, error)
}
