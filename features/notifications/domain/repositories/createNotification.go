package notificationdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type CreateNotificationTransactionComponent interface {
	PreTransaction(
		createNotificationInput *model.InternalCreateNotification,
	) (*model.InternalCreateNotification, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createNotificationInput *model.InternalCreateNotification,
	) (*model.Notification, error)
}
