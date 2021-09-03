package moudomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ApproveUpdateMouTransactionComponent interface {
	PreTransaction(
		updateMouInput *model.InternalUpdateMou,
	) (*model.InternalUpdateMou, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateMouInput *model.InternalUpdateMou,
	) (*model.Mou, error)
}

type ApproveUpdateMouRepository interface {
	RunTransaction(
		updateMouInput *model.InternalUpdateMou,
	) (*model.Mou, error)
}
