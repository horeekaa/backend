package taggingdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type UpdateTaggingUsecaseComponent interface {
	Validation(
		updateTaggingInput *model.InternalUpdateTagging,
	) (*model.InternalUpdateTagging, error)
}

type UpdateTaggingTransactionComponent interface {
	SetValidation(usecaseComponent UpdateTaggingUsecaseComponent) (bool, error)

	PreTransaction(
		updateTaggingInput *model.InternalUpdateTagging,
	) (*model.InternalUpdateTagging, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateTaggingInput *model.InternalUpdateTagging,
	) ([]*model.Tagging, error)
}

type UpdateTaggingRepository interface {
	SetValidation(usecaseComponent UpdateTaggingUsecaseComponent) (bool, error)
	RunTransaction(
		updateTaggingInput *model.InternalUpdateTagging,
	) ([]*model.Tagging, error)
}
