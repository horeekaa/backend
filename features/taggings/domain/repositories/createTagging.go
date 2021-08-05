package taggingdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type CreateTaggingUsecaseComponent interface {
	Validation(
		createTaggingInput *model.InternalCreateTagging,
	) (*model.InternalCreateTagging, error)
}

type CreateTaggingTransactionComponent interface {
	SetValidation(usecaseComponent CreateTaggingUsecaseComponent) (bool, error)

	PreTransaction(
		createTaggingInput *model.InternalCreateTagging,
	) (*model.InternalCreateTagging, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createTaggingInput *model.InternalCreateTagging,
	) ([]*model.Tagging, error)
}

type CreateTaggingRepository interface {
	SetValidation(usecaseComponent CreateTaggingUsecaseComponent) (bool, error)
	RunTransaction(
		createTaggingInput *model.InternalCreateTagging,
	) ([]*model.Tagging, error)
}
