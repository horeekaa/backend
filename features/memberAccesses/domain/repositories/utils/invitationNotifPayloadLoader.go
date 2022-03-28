package memberaccessdomainrepositoryutilityinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type InvitationPayloadLoader interface {
	Execute(
		notification *model.InternalCreateNotification,
	) (bool, error)
}
