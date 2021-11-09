package mongodbpurchaseordertosupplydatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PurchaseOrderToSupplyDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderToSupply, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderToSupply, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.PurchaseOrderToSupply, error)
	Create(input *model.DatabaseCreatePurchaseOrderToSupply, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderToSupply, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdatePurchaseOrderToSupply, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderToSupply, error)
}
