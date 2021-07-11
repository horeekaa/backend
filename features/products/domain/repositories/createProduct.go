package productdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type CreateProductUsecaseComponent interface {
	Validation(
		createProductInput *model.InternalCreateProduct,
	) (*model.InternalCreateProduct, error)
}

type CreateProductTransactionComponent interface {
	SetValidation(usecaseComponent CreateProductUsecaseComponent) (bool, error)

	PreTransaction(
		createProductInput *model.InternalCreateProduct,
	) (*model.InternalCreateProduct, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createProductInput *model.InternalCreateProduct,
	) (*model.Product, error)
}

type CreateProductRepository interface {
	SetValidation(usecaseComponent CreateProductUsecaseComponent) (bool, error)
	RunTransaction(
		createProductInput *model.InternalCreateProduct,
	) (*model.Product, error)
}
