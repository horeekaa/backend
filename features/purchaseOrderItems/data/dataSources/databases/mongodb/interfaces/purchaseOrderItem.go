package mongodbpurchaseorderitemdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PurchaseOrderItemDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderItem, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderItem, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.PurchaseOrderItem, error)
	Create(input *model.DatabaseCreatePurchaseOrderItem, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderItem, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdatePurchaseOrderItem, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderItem, error)
}
