package mongodbnotificationdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
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
	_, err := notificationDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.DatabaseNotification
	_, err = notificationDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (notificationDataSourceMongo *notificationDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdateNotification,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.DatabaseNotification, error) {
	_, err := notificationDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.DatabaseNotification
	_, err = notificationDataSourceMongo.basicOperation.Update(
		updateCriteria,
		map[string]interface{}{
			"$set": updateData,
		},
		&output,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (notificationDataSourceMongo *notificationDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdateNotification,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()
	existingObject, err := notificationDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			"/notificationDataSource/update",
			nil,
		)
	}

	input.UpdatedAt = &currentTime

	return true, nil
}

func (notificationDataSourceMongo *notificationDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreateNotification,
) (bool, error) {
	currentTime := time.Now()

	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime

	return true, nil
}
