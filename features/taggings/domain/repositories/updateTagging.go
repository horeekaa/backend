package taggingdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type UpdateTaggingUsecaseComponent interface {
	Validation(
		updateTaggingInput *model.InternalBulkUpdateTagging,
	) (*model.InternalBulkUpdateTagging, error)
}

type UpdateTaggingTransactionComponent interface {
	SetValidation(usecaseComponent UpdateTaggingUsecaseComponent) (bool, error)

	PreTransaction(
		updateTaggingInput *model.InternalBulkUpdateTagging,
	) (*model.InternalBulkUpdateTagging, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateTaggingInput *model.InternalBulkUpdateTagging,
	) ([]*model.Tagging, error)
}

type UpdateTaggingRepository interface {
	SetValidation(usecaseComponent UpdateTaggingUsecaseComponent) (bool, error)
	RunTransaction(
		updateTaggingInput *model.InternalBulkUpdateTagging,
	) ([]*model.Tagging, error)
}
