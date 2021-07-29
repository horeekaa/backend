package productvariantdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type UpdateProductVariantTransactionComponent interface {
	PreTransaction(
		updateProductVariantInput *model.InternalUpdateProductVariant,
	) (*model.InternalUpdateProductVariant, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateProductVariantInput *model.InternalUpdateProductVariant,
	) (*model.ProductVariant, error)
}
