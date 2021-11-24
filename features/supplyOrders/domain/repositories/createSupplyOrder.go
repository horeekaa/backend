package supplyorderdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateSupplyOrderTransactionComponent interface {
	PreTransaction(
		createSupplyOrderInput *model.InternalCreateSupplyOrder,
	) (*model.InternalCreateSupplyOrder, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createSupplyOrderInput *model.InternalCreateSupplyOrder,
	) (*model.SupplyOrder, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}

type CreateSupplyOrderRepository interface {
	RunTransaction(
		createSupplyOrderInput *model.InternalCreateSupplyOrder,
	) ([]*model.SupplyOrder, error)
}
