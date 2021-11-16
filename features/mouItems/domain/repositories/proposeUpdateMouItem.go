package mouitemdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ProposeUpdateMouItemTransactionComponent interface {
	PreTransaction(
		updateMouItemInput *model.InternalUpdateMouItem,
	) (*model.InternalUpdateMouItem, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateMouItemInput *model.InternalUpdateMouItem,
	) (*model.MouItem, error)
}
