package productdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateProductRepository struct {
	approveUpdateProductTransactionComponent productdomainrepositoryinterfaces.ApproveUpdateProductTransactionComponent
	mongoDBTransaction                       mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateProductRepository(
	approveUpdateProductRepositoryTransactionComponent productdomainrepositoryinterfaces.ApproveUpdateProductTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (productdomainrepositoryinterfaces.ApproveUpdateProductRepository, error) {
	approveUpdateProductRepo := &approveUpdateProductRepository{
		approveUpdateProductRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		approveUpdateProductRepo,
		"ApproveUpdateProductRepository",
	)

	return approveUpdateProductRepo, nil
}

func (updateOrgRepo *approveUpdateProductRepository) SetValidation(
	usecaseComponent productdomainrepositoryinterfaces.ApproveUpdateProductUsecaseComponent,
) (bool, error) {
	updateOrgRepo.approveUpdateProductTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateOrgRepo *approveUpdateProductRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.approveUpdateProductTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateProduct),
	)
}

func (updateOrgRepo *approveUpdateProductRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.approveUpdateProductTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateProduct),
	)
}

func (updateOrgRepo *approveUpdateProductRepository) RunTransaction(
	input *model.InternalUpdateProduct,
) (*model.Product, error) {
	output, err := updateOrgRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Product), err
}
