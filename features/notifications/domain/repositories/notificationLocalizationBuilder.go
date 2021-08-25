package notificationdomainrepositoryinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type NotificationLocalizationBuilder interface {
	Execute(
		input *model.DatabaseNotification,
		output *model.Notification,
	) (bool, error)
}
