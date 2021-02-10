package mongooperationinterfaces

import (
	mongooperationmodels "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type BasicOperation interface {
	FindByID(ID interface{}, operationOptions *mongooperationmodels.OperationOptions) (*mongo.SingleResult, error)
	FindOne(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) (*mongo.SingleResult, error)
	Find(query map[string]interface{}, cursorDecoder func(cursorObject *mongooperationmodels.CursorObject) (interface{}, error), operationOptions *mongooperationmodels.OperationOptions) (*bool, error)
	Create(input interface{}, operationOptions *mongooperationmodels.OperationOptions) (*mongooperationmodels.CreateOperationOutput, error)
	Update(ID interface{}, updateData interface{}, operationOptions *mongooperationmodels.OperationOptions) (*mongo.SingleResult, error)
}
