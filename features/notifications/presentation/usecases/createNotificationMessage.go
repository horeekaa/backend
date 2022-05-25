package notificationpresentationusecaseinterfaces

import (
	notificationpresentationusecasetypes "github.com/horeekaa/backend/features/notifications/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreateNotificationMessage interface {
	Execute(input notificationpresentationusecasetypes.CreateNotificationMessageUsecaseInput) (*model.NotificationMessage, error)
}
