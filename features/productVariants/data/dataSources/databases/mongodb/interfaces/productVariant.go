package mongodbproductvariantdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductVariantDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.ProductVariant, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.ProductVariant, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.ProductVariant, error)
	Create(input *model.DatabaseCreateProductVariant, operationOptions *mongodbcoretypes.OperationOptions) (*model.ProductVariant, error)
	Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateProductVariant, operationOptions *mongodbcoretypes.OperationOptions) (*model.ProductVariant, error)
	GenerateObjectID() primitive.ObjectID
}
