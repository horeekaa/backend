package mongodbmoudatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MouDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Mou, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Mou, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Mou, error)
	Create(input *model.DatabaseCreateMou, operationOptions *mongodbcoretypes.OperationOptions) (*model.Mou, error)
	Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateMou, operationOptions *mongodbcoretypes.OperationOptions) (*model.Mou, error)
	GenerateObjectID() primitive.ObjectID
}
