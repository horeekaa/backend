package mouitemdomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type AgreedProductLoader interface {
	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		mouItem *model.DatabaseCreateMouItem,
	) (bool, error)
}
