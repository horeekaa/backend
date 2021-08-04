package mongodbtaggingdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaggingDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tagging, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tagging, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Tagging, error)
	Create(input *model.DatabaseCreateTagging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tagging, error)
	Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateTagging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tagging, error)
}
