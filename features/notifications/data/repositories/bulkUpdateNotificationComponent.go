package notificationdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasenotificationdatasourceinterfaces "github.com/horeekaa/backend/features/notifications/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type bulkUpdateNotificationTransactionComponent struct {
	notificationDataSource   databasenotificationdatasourceinterfaces.NotificationDataSource
	notifLocalizationBuilder notificationdomainrepositoryinterfaces.NotificationLocalizationBuilder
}

func NewBulkUpdateNotificationTransactionComponent(
	notificationDataSource databasenotificationdatasourceinterfaces.NotificationDataSource,
	notifLocalizationBuilder notificationdomainrepositoryinterfaces.NotificationLocalizationBuilder,
) (notificationdomainrepositoryinterfaces.BulkUpdateNotificationTransactionComponent, error) {
	return &bulkUpdateNotificationTransactionComponent{
		notificationDataSource:   notificationDataSource,
		notifLocalizationBuilder: notifLocalizationBuilder,
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

		databaseNotification, err := bulkUpdateNotificationComp.notificationDataSource.GetMongoDataSource().Update(
			id,
			notificationToUpdate,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/bulkProposeUpdateNotification",
				err,
			)
		}

		notification := &model.Notification{}
		jsonOutput, _ := json.Marshal(databaseNotification)
		json.Unmarshal(jsonOutput, notification)

		bulkUpdateNotificationComp.notifLocalizationBuilder.Execute(databaseNotification, notification)

		notifications = append(notifications, notification)
	}
	return notifications, nil
}
