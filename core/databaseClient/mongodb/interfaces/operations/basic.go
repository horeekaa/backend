package mongodbcoreoperationinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BasicOperation interface {
	SetCollection(collectionName string) bool
	GetCollectionName() string
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error)
	Find(
		query map[string]interface{},
		paginationOpt *mongodbcoretypes.PaginationOptions,
		cursorDecoder func(cursorObject *mongo.Cursor) (interface{}, error),
		operationOptions *mongodbcoretypes.OperationOptions,
	) (*bool, error)
	Create(input interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongodbcoretypes.CreateOperationOutput, error)
	Update(ID primitive.ObjectID, updateData interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*mongo.SingleResult, error)
}
