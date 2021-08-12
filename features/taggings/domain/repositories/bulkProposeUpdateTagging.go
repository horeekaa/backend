package taggingdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type BulkProposeUpdateTaggingUsecaseComponent interface {
	Validation(
		bulkProposeUpdateTaggingInput *model.InternalBulkUpdateTagging,
	) (*model.InternalBulkUpdateTagging, error)
}

type BulkProposeUpdateTaggingTransactionComponent interface {
	SetValidation(usecaseComponent BulkProposeUpdateTaggingUsecaseComponent) (bool, error)

	PreTransaction(
		bulkProposeUpdateTaggingInput *model.InternalBulkUpdateTagging,
	) (*model.InternalBulkUpdateTagging, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		bulkProposeUpdateTaggingInput *model.InternalBulkUpdateTagging,
	) ([]*model.Tagging, error)
}

type BulkProposeUpdateTaggingRepository interface {
	SetValidation(usecaseComponent BulkProposeUpdateTaggingUsecaseComponent) (bool, error)
	RunTransaction(
		bulkProposeUpdateTaggingInput *model.InternalBulkUpdateTagging,
	) ([]*model.Tagging, error)
}
