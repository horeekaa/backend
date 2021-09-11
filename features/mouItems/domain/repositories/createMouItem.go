package mouitemdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateMouItemTransactionComponent interface {
	PreTransaction(
		createMouItemInput *model.InternalCreateMouItem,
	) (*model.InternalCreateMouItem, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createMouItemInput *model.InternalCreateMouItem,
	) (*model.MouItem, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}
