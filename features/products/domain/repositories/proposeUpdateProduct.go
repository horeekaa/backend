package productdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ProposeUpdateProductUsecaseComponent interface {
	Validation(
		updateProductInput *model.InternalUpdateProduct,
	) (*model.InternalUpdateProduct, error)
}

type ProposeUpdateProductTransactionComponent interface {
	SetValidation(usecaseComponent ProposeUpdateProductUsecaseComponent) (bool, error)

	PreTransaction(
		updateProductInput *model.InternalUpdateProduct,
	) (*model.InternalUpdateProduct, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateProductInput *model.InternalUpdateProduct,
	) (*model.Product, error)
}

type ProposeUpdateProductRepository interface {
	SetValidation(usecaseComponent ProposeUpdateProductUsecaseComponent) (bool, error)
	RunTransaction(
		updateProductInput *model.InternalUpdateProduct,
	) (*model.Product, error)
}
