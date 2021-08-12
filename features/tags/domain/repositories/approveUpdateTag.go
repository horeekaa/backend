package tagdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ApproveUpdateTagUsecaseComponent interface {
	Validation(
		updateTagInput *model.InternalUpdateTag,
	) (*model.InternalUpdateTag, error)
}

type ApproveUpdateTagTransactionComponent interface {
	SetValidation(usecaseComponent ApproveUpdateTagUsecaseComponent) (bool, error)

	PreTransaction(
		updateTagInput *model.InternalUpdateTag,
	) (*model.InternalUpdateTag, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateTagInput *model.InternalUpdateTag,
	) (*model.Tag, error)
}

type ApproveUpdateTagRepository interface {
	RunTransaction(
		updateTagInput *model.InternalUpdateTag,
	) (*model.Tag, error)
}
