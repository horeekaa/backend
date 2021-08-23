package notificationpresentationusecaseinterfaces

import (
	notificationpresentationusecasetypes "github.com/horeekaa/backend/features/notifications/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type BulkUpdateNotificationUsecase interface {
	Execute(input notificationpresentationusecasetypes.BulkUpdateNotificationUsecaseInput) ([]*model.Notification, error)
}
