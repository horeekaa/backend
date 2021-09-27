package purchaseorderdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePurchaseOrderTransactionComponent interface {
	PreTransaction(
		createPurchaseOrderInput *model.InternalCreatePurchaseOrder,
	) (*model.InternalCreatePurchaseOrder, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createPurchaseOrderInput *model.InternalCreatePurchaseOrder,
	) (*model.PurchaseOrder, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}

type CreatePurchaseOrderRepository interface {
	RunTransaction(
		createPurchaseOrderInput *model.InternalCreatePurchaseOrder,
	) ([]*model.PurchaseOrder, error)
}
