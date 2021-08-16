package mongodbnotificationdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	mongodbnotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type notificationDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewNotificationDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbnotificationdatasourceinterfaces.NotificationDataSourceMongo, error) {
	basicOperation.SetCollection("notifications")
	return &notificationDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (notificationDataSourceMongo *notificationDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (notificationDataSourceMongo *notificationDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.DatabaseNotification, error) {
	var output model.DatabaseNotification
	_, err := notificationDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (notificationDataSourceMongo *notificationDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.DatabaseNotification, error) {
	var output model.DatabaseNotification
	_, err := notificationDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (notificationDataSourceMongo *notificationDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.DatabaseNotification, error) {
	var notifications = []*model.DatabaseNotification{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var notification model.DatabaseNotification
		if err := cursor.Decode(&notification); err != nil {
			return err
		}
		notifications = append(notifications, &notification)
		return nil
	}
	_, err := notificationDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return notifications, err
}

func (notificationDataSourceMongo *notificationDataSourceMongo) Create(input *model.DatabaseCreateNotification, operationOptions *mongodbcoretypes.OperationOptions) (*model.DatabaseNotification, error) {
	defaultedInput, err := notificationDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.DatabaseNotification
	_, err = notificationDataSourceMongo.basicOperation.Create(*defaultedInput.CreateNotification, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (notificationDataSourceMongo *notificationDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateNotification, operationOptions *mongodbcoretypes.OperationOptions) (*model.DatabaseNotification, error) {
	updateData.ID = ID
	defaultedInput, err := notificationDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.DatabaseNotification
	_, err = notificationDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateNotification, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type setnotificationDefaultValuesOutput struct {
	CreateNotification *model.DatabaseCreateNotification
	UpdateNotification *model.DatabaseUpdateNotification
}

func (notificationDataSourceMongo *notificationDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setnotificationDefaultValuesOutput, error) {
	currentTime := time.Now()

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.DatabaseUpdateNotification)
		_, err := notificationDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setnotificationDefaultValuesOutput{
			UpdateNotification: &updateInput,
		}, nil
	}
	createInput := (input).(model.DatabaseCreateNotification)
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime

	return &setnotificationDefaultValuesOutput{
		CreateNotification: &createInput,
	}, nil
}
