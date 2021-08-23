package notificationpresentationusecaseinterfaces

import (
	notificationpresentationusecasetypes "github.com/horeekaa/backend/features/notifications/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllNotificationUsecase interface {
	Execute(input notificationpresentationusecasetypes.GetAllNotificationUsecaseInput) ([]*model.Notification, error)
}
