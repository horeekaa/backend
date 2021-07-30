package productvariantdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateProductVariantTransactionComponent interface {
	PreTransaction(
		createProductVariantInput *model.InternalCreateProductVariant,
	) (*model.InternalCreateProductVariant, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createProductVariantInput *model.InternalCreateProductVariant,
	) (*model.ProductVariant, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}
