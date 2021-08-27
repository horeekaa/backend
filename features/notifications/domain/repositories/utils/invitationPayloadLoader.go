package notificationdomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type InvitationPayloadLoader interface {
	TransactionBody(
		operationOptions *mongodbcoretypes.OperationOptions,
		notification *model.InternalCreateNotification,
	) (bool, error)
}
