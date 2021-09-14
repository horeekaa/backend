package mongodbmouitemdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MouItemDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.MouItem, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MouItem, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.MouItem, error)
	Create(input *model.DatabaseCreateMouItem, operationOptions *mongodbcoretypes.OperationOptions) (*model.MouItem, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateMouItem, operationOptions *mongodbcoretypes.OperationOptions) (*model.MouItem, error)
	GenerateObjectID() primitive.ObjectID
}
