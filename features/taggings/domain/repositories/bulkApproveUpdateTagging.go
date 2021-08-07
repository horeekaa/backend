package taggingdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type BulkApproveUpdateTaggingUsecaseComponent interface {
	Validation(
		bulkUpdateTaggingInput *model.InternalBulkUpdateTagging,
	) (*model.InternalBulkUpdateTagging, error)
}

type BulkApproveUpdateTaggingTransactionComponent interface {
	SetValidation(usecaseComponent BulkApproveUpdateTaggingUsecaseComponent) (bool, error)

	PreTransaction(
		bulkUpdateTaggingInput *model.InternalBulkUpdateTagging,
	) (*model.InternalBulkUpdateTagging, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		bulkUpdateTaggingInput *model.InternalBulkUpdateTagging,
	) ([]*model.Tagging, error)
}

type BulkApproveUpdateTaggingRepository interface {
	RunTransaction(
		bulkUpdateTaggingInput *model.InternalBulkUpdateTagging,
	) ([]*model.Tagging, error)
}
