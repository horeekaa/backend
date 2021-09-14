package mongodbnotificationdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.DatabaseNotification, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.DatabaseNotification, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.DatabaseNotification, error)
	Create(input *model.DatabaseCreateNotification, operationOptions *mongodbcoretypes.OperationOptions) (*model.DatabaseNotification, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateNotification, operationOptions *mongodbcoretypes.OperationOptions) (*model.DatabaseNotification, error)
	GenerateObjectID() primitive.ObjectID
}
