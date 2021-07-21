package mongodbproductdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Product, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Product, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Product, error)
	Create(input *model.DatabaseCreateProduct, operationOptions *mongodbcoretypes.OperationOptions) (*model.Product, error)
	Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateProduct, operationOptions *mongodbcoretypes.OperationOptions) (*model.Product, error)
	GenerateObjectID() primitive.ObjectID
}
