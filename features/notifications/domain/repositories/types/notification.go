package notificationdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllNotificationInput struct {
	FilterFields  *model.NotificationFilterFields
	PaginationOpt *model.PaginationOptionInput
}
