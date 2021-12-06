package mongodbsupplyorderitemdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SupplyOrderItemDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrderItem, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrderItem, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.SupplyOrderItem, error)
	Create(input *model.DatabaseCreateSupplyOrderItem, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrderItem, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateSupplyOrderItem, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrderItem, error)
	GenerateObjectID() primitive.ObjectID
}
