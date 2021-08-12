package taggingdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BulkCreateTaggingUsecaseComponent interface {
	Validation(
		bulkCreateTaggingInput *model.InternalCreateTagging,
	) (*model.InternalCreateTagging, error)
}

type BulkCreateTaggingTransactionComponent interface {
	SetValidation(usecaseComponent BulkCreateTaggingUsecaseComponent) (bool, error)

	PreTransaction(
		bulkCreateTaggingInput *model.InternalCreateTagging,
	) (*model.InternalCreateTagging, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		bulkCreateTaggingInput *model.InternalCreateTagging,
	) ([]*model.Tagging, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}

type BulkCreateTaggingRepository interface {
	SetValidation(usecaseComponent BulkCreateTaggingUsecaseComponent) (bool, error)
	RunTransaction(
		bulkCreateTaggingInput *model.InternalCreateTagging,
	) ([]*model.Tagging, error)
}
