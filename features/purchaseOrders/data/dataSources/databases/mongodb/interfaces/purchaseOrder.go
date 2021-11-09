package mongodbpurchaseorderdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PurchaseOrderDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrder, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrder, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.PurchaseOrder, error)
	Create(input *model.DatabaseCreatePurchaseOrder, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrder, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdatePurchaseOrder, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrder, error)
	GenerateObjectID() primitive.ObjectID
}
