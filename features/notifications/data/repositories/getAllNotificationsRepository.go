package notificationdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	notificationdomainrepositorytypes "github.com/horeekaa/backend/features/notifications/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getAllNotificationRepository struct {
	notificationDataSource   databasenotificationdatasourceinterfaces.NotificationDataSource
	notifLocalizationBuilder notificationdomainrepositoryinterfaces.NotificationLocalizationBuilder
}

func NewGetAllNotificationRepository(
	notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
	notifLocalizationBuilder notificationdomainrepositoryinterfaces.NotificationLocalizationBuilder,
) (notificationdomainrepositoryinterfaces.GetAllNotificationRepository, error) {
	return &getAllNotificationRepository{
		notificationDataSource,
		notifLocalizationBuilder,
	}, nil
}

func (getAllNotificationRepo *getAllNotificationRepository) Execute(
	input notificationdomainrepositorytypes.GetAllNotificationInput,
) ([]*model.Notification, error) {
	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(input.FilterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	databaseNotifs, err := getAllNotificationRepo.notificationDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAllNotification",
			err,
		)
	}

	notifications := []*model.Notification{}
	for _, notif := range databaseNotifs {
		notification := &model.Notification{}
		jsonTemp, _ := json.Marshal(notif)
		json.Unmarshal(jsonTemp, notification)

		getAllNotificationRepo.notifLocalizationBuilder.Execute(notif, notification)

		notifications = append(notifications, notification)
	}

	return notifications, nil
}
