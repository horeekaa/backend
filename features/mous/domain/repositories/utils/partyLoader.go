package moudomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type PartyLoader interface {
	TransactionBody(
		operationOptions *mongodbcoretypes.OperationOptions,
		input *model.PartyInput,
		output *model.InternalPartyInput,
	) (bool, error)
}
