package notificationdomainrepositoryinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type NotificationLocalizationBuilder interface {
	Execute(
		input *model.InternalCreateNotification,
		output *model.Notification,
	) (bool, error)
}
