package notificationdomainrepositoryinterfaces

import (
	notificationdomainrepositorytypes "github.com/horeekaa/backend/features/notifications/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllNotificationRepository interface {
	Execute(filterFields notificationdomainrepositorytypes.GetAllNotificationInput) ([]*model.Notification, error)
}
