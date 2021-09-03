package moudomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateMouTransactionComponent interface {
	PreTransaction(
		createMouInput *model.InternalCreateMou,
	) (*model.InternalCreateMou, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createMouInput *model.InternalCreateMou,
	) (*model.Mou, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}

type CreateMouRepository interface {
	RunTransaction(
		createMouInput *model.InternalCreateMou,
	) (*model.Mou, error)
}
