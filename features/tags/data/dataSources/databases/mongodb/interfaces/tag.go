package mongodbtagdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tag, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tag, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Tag, error)
	Create(input *model.DatabaseCreateTag, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tag, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateTag, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tag, error)
	GenerateObjectID() primitive.ObjectID
}
