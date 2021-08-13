package mongodbnotificationdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Notification, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Notification, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Notification, error)
	Create(input *model.DatabaseCreateNotification, operationOptions *mongodbcoretypes.OperationOptions) (*model.Notification, error)
	Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateNotification, operationOptions *mongodbcoretypes.OperationOptions) (*model.Notification, error)
	GenerateObjectID() primitive.ObjectID
}
