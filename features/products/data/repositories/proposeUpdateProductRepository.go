package productdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateProductRepository struct {
	proposeUpdateProductTransactionComponent productdomainrepositoryinterfaces.ProposeUpdateProductTransactionComponent
	mongoDBTransaction                       mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateProductRepository(
	proposeUpdateProductRepositoryTransactionComponent productdomainrepositoryinterfaces.ProposeUpdateProductTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (productdomainrepositoryinterfaces.ProposeUpdateProductRepository, error) {
	proposeUpdateProductRepo := &proposeUpdateProductRepository{
		proposeUpdateProductRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateProductRepo,
		"ProposeUpdateProductRepository",
	)

	return proposeUpdateProductRepo, nil
}

func (updateOrgRepo *proposeUpdateProductRepository) SetValidation(
	usecaseComponent productdomainrepositoryinterfaces.ProposeUpdateProductUsecaseComponent,
) (bool, error) {
	updateOrgRepo.proposeUpdateProductTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateOrgRepo *proposeUpdateProductRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.proposeUpdateProductTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateProduct),
	)
}

func (updateOrgRepo *proposeUpdateProductRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.proposeUpdateProductTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateProduct),
	)
}

func (updateOrgRepo *proposeUpdateProductRepository) RunTransaction(
	input *model.InternalUpdateProduct,
) (*model.Product, error) {
	output, err := updateOrgRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Product), err
}
