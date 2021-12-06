package mongodbsupplyorderdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SupplyOrderDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrder, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrder, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.SupplyOrder, error)
	Create(input *model.DatabaseCreateSupplyOrder, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrder, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateSupplyOrder, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrder, error)
	GenerateObjectID() primitive.ObjectID
}
