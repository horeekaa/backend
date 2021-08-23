package notificationdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type bulkUpdateNotificationRepository struct {
	bulkUpdateNotificationTransactionComponent notificationdomainrepositoryinterfaces.BulkUpdateNotificationTransactionComponent
	mongoDBTransaction                         mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewBulkUpdateNotificationRepository(
	bulkUpdateNotificationRepositoryTransactionComponent notificationdomainrepositoryinterfaces.BulkUpdateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (notificationdomainrepositoryinterfaces.BulkUpdateNotificationRepository, error) {
	bulkUpdateNotificationRepo := &bulkUpdateNotificationRepository{
		bulkUpdateNotificationRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		bulkUpdateNotificationRepo,
		"BulkUpdateNotificationRepository",
	)

	return bulkUpdateNotificationRepo, nil
}

func (bulkUpdatenotificationRepo *bulkUpdateNotificationRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return bulkUpdatenotificationRepo.bulkUpdateNotificationTransactionComponent.PreTransaction(
		input.(*model.InternalBulkUpdateNotification),
	)
}

func (bulkUpdatenotificationRepo *bulkUpdateNotificationRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return bulkUpdatenotificationRepo.bulkUpdateNotificationTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalBulkUpdateNotification),
	)
}

func (bulkUpdatenotificationRepo *bulkUpdateNotificationRepository) RunTransaction(
	input *model.InternalBulkUpdateNotification,
) ([]*model.Notification, error) {
	output, err := bulkUpdatenotificationRepo.mongoDBTransaction.RunTransaction(input)
	return (output).([]*model.Notification), err
}
