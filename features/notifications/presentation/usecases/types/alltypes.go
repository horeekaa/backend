package notificationpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type GetAllNotificationUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.NotificationFilterFields
	PaginationOps *model.PaginationOptionInput
}

type BulkUpdateNotificationUsecaseInput struct {
	Context                context.Context
	BulkUpdateNotification *model.BulkUpdateNotification
}

type CreateNotificationMessageUsecaseInput struct {
	Context      context.Context
	Notification *model.Notification
}
