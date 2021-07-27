package productvariantdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type CreateProductVariantTransactionComponent interface {
	PreTransaction(
		createProductVariantInput *model.InternalCreateProductVariant,
	) (*model.InternalCreateProductVariant, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createProductVariantInput *model.InternalCreateProductVariant,
	) (*model.ProductVariant, error)
}
