package tagdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ProposeUpdateTagUsecaseComponent interface {
	Validation(
		updateTagInput *model.InternalUpdateTag,
	) (*model.InternalUpdateTag, error)
}

type ProposeUpdateTagTransactionComponent interface {
	SetValidation(usecaseComponent ProposeUpdateTagUsecaseComponent) (bool, error)

	PreTransaction(
		updateTagInput *model.InternalUpdateTag,
	) (*model.InternalUpdateTag, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateTagInput *model.InternalUpdateTag,
	) (*model.Tag, error)
}

type ProposeUpdateTagRepository interface {
	SetValidation(usecaseComponent ProposeUpdateTagUsecaseComponent) (bool, error)
	RunTransaction(
		updateTagInput *model.InternalUpdateTag,
	) (*model.Tag, error)
}
