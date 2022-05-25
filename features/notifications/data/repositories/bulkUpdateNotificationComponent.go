package notificationdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type bulkUpdateNotificationTransactionComponent struct {
	notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource
	pathIdentity           string
}

func NewBulkUpdateNotificationTransactionComponent(
	notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
) (notificationdomainrepositoryinterfaces.BulkUpdateNotificationTransactionComponent, error) {
	return &bulkUpdateNotificationTransactionComponent{
		notificationDataSource: notificationDataSource,
		pathIdentity:           "BulkUpdateNotificationComponent",
	}, nil
}

func (bulkUpdateNotificationComp *bulkUpdateNotificationTransactionComponent) PreTransaction(
	input *model.InternalBulkUpdateNotification,
) (*model.InternalBulkUpdateNotification, error) {
	return input, nil
}

func (bulkUpdateNotificationComp *bulkUpdateNotificationTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalBulkUpdateNotification,
) ([]*model.Notification, error) {
	notifications := []*model.Notification{}
	jsonTemp, _ := json.Marshal(input)
	for _, id := range input.IDs {
		notificationToUpdate := &model.DatabaseUpdateNotification{}
		json.Unmarshal(jsonTemp, notificationToUpdate)

		currentTime := time.Now().UTC()
		notificationToUpdate.UpdatedAt = &currentTime

		databaseNotification, err := bulkUpdateNotificationComp.notificationDataSource.GetMongoDataSource().Update(
			map[string]interface{}{
				"_id": id,
			},
			notificationToUpdate,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				bulkUpdateNotificationComp.pathIdentity,
				err,
			)
		}

		notification := &model.Notification{}
		jsonOutput, _ := json.Marshal(databaseNotification)
		json.Unmarshal(jsonOutput, notification)

		notifications = append(notifications, notification)
	}
	return notifications, nil
}
