package notificationdomainrepositories

import (
	"encoding/json"

	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	notificationdomainrepositorytypes "github.com/horeekaa/backend/features/notifications/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllNotificationRepository struct {
	notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource
	mongoQueryBuilder      mongodbcorequerybuilderinterfaces.MongoQueryBuilder
	pathIdentity           string
}

func NewGetAllNotificationRepository(
	notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (notificationdomainrepositoryinterfaces.GetAllNotificationRepository, error) {
	return &getAllNotificationRepository{
		notificationDataSource,
		mongoQueryBuilder,
		"GetAllNotificationRepository",
	}, nil
}

func (getAllNotificationRepo *getAllNotificationRepository) Execute(
	input notificationdomainrepositorytypes.GetAllNotificationInput,
) ([]*model.Notification, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllNotificationRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	databaseNotifs, err := getAllNotificationRepo.notificationDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAllNotificationRepo.pathIdentity,
			err,
		)
	}

	notifications := []*model.Notification{}
	for _, notif := range databaseNotifs {
		notification := &model.Notification{}
		jsonTemp, _ := json.Marshal(notif)
		json.Unmarshal(jsonTemp, notification)

		notifications = append(notifications, notification)
	}

	return notifications, nil
}
