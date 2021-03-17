package mongodbcoreoperationinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type BasicOperation interface {
	FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error)
	Find(query map[string]interface{}, cursorDecoder func(cursorObject *mongodbcoretypes.CursorObject) (interface{}, error), operationOptions *mongodbcoretypes.OperationOptions) (*bool, error)
	Create(input interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongodbcoretypes.CreateOperationOutput, error)
	Update(ID interface{}, updateData interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error)
}
