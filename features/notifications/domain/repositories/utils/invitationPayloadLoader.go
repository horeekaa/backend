package notificationdomainrepositoryutilityinterfaces

import "github.com/horeekaa/backend/model"

type InvitationPayloadLoader interface {
	Execute(
		notification *model.DatabaseNotification,
	) (bool, error)
}
