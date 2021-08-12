package tagdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTagUsecaseComponent interface {
	Validation(
		createTagInput *model.InternalCreateTag,
	) (*model.InternalCreateTag, error)
}

type CreateTagTransactionComponent interface {
	SetValidation(usecaseComponent CreateTagUsecaseComponent) (bool, error)

	PreTransaction(
		createTagInput *model.InternalCreateTag,
	) (*model.InternalCreateTag, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createTagInput *model.InternalCreateTag,
	) (*model.Tag, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}

type CreateTagRepository interface {
	RunTransaction(
		createTagInput *model.InternalCreateTag,
	) (*model.Tag, error)
}
