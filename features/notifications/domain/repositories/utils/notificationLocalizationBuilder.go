package notificationdomainrepositoryutilityinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type NotificationLocalizationBuilder interface {
	Execute(
		input *model.DatabaseNotification,
		output *model.Notification,
		language string,
	) (bool, error)
}
