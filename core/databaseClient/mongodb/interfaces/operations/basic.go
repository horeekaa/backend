package mongodbcoreoperationinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BasicOperation interface {
	SetCollection(collectionName string) bool
	GetCollectionName() string
	FindByID(ID primitive.ObjectID, output interface{}, operationOptions *mongodbcoretypes.OperationOptions) (bool, error)
	FindOne(query map[string]interface{}, output interface{}, operationOptions *mongodbcoretypes.OperationOptions) (bool, error)
	Find(
		query map[string]interface{},
		paginationOpt *mongodbcoretypes.PaginationOptions,
		output interface{},
		operationOptions *mongodbcoretypes.OperationOptions,
	) (bool, error)
	Create(input interface{}, output interface{}, operationOptions *mongodbcoretypes.OperationOptions) (bool, error)
	Update(ID primitive.ObjectID, updateData interface{}, output interface{}, operationOptions *mongodbcoretypes.OperationOptions) (bool, error)
}
