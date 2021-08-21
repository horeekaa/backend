package notificationdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type BulkUpdateNotificationTransactionComponent interface {
	PreTransaction(
		updateNotificationInput *model.InternalBulkUpdateNotification,
	) (*model.InternalBulkUpdateNotification, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateNotificationInput *model.InternalBulkUpdateNotification,
	) ([]*model.Notification, error)
}
