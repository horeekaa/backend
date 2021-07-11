package productdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createProductRepository struct {
	createProductTransactionComponent productdomainrepositoryinterfaces.CreateProductTransactionComponent
	mongoDBTransaction                mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateProductRepository(
	createProductRepositoryTransactionComponent productdomainrepositoryinterfaces.CreateProductTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (productdomainrepositoryinterfaces.CreateProductRepository, error) {
	createProductRepo := &createProductRepository{
		createProductRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createProductRepo,
		"CreateProductRepository",
	)

	return createProductRepo, nil
}

func (createProdRepo *createProductRepository) SetValidation(
	usecaseComponent productdomainrepositoryinterfaces.CreateProductUsecaseComponent,
) (bool, error) {
	createProdRepo.createProductTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (createProdRepo *createProductRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createProdRepo.createProductTransactionComponent.PreTransaction(
		input.(*model.InternalCreateProduct),
	)
}

func (createProdRepo *createProductRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return createProdRepo.createProductTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalCreateProduct),
	)
}

func (createProdRepo *createProductRepository) RunTransaction(
	input *model.InternalCreateProduct,
) (*model.Product, error) {
	output, err := createProdRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Product), err
}
