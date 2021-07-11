package productdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ApproveUpdateProductUsecaseComponent interface {
	Validation(
		updateProductInput *model.InternalUpdateProduct,
	) (*model.InternalUpdateProduct, error)
}

type ApproveUpdateProductTransactionComponent interface {
	SetValidation(usecaseComponent ApproveUpdateProductUsecaseComponent) (bool, error)

	PreTransaction(
		updateProductInput *model.InternalUpdateProduct,
	) (*model.InternalUpdateProduct, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateProductInput *model.InternalUpdateProduct,
	) (*model.Product, error)
}

type ApproveUpdateProductRepository interface {
	SetValidation(usecaseComponent ApproveUpdateProductUsecaseComponent) (bool, error)
	RunTransaction(
		updateProductInput *model.InternalUpdateProduct,
	) (*model.Product, error)
}
