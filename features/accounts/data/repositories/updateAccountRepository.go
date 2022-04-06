package accountdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateAccountRepository struct {
	updateAccountTransactionComponent accountdomainrepositoryinterfaces.UpdateAccountTransactionComponent
	mongoDBTransaction                       mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewUpdateAccountRepository(
	updateAccountRepositoryTransactionComponent accountdomainrepositoryinterfaces.UpdateAccountTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (accountdomainrepositoryinterfaces.UpdateAccountRepository, error) {
	updateAccountRepo := &updateAccountRepository{
		updateAccountRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		updateAccountRepo,
		"UpdateAccountRepository",
	)

	return updateAccountRepo, nil
}

func (updateAccountRepo *updateAccountRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateAccountRepo.updateAccountTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateAccount),
	)
}

func (updateAccountRepo *updateAccountRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	accountToUpdate := input.(*model.InternalUpdateAccount)

	return updateAccountRepo.updateAccountTransactionComponent.TransactionBody(
		operationOption,
		accountToUpdate,
	)
}

func (updateAccountRepo *updateAccountRepository) RunTransaction(
	input *model.InternalUpdateAccount,
) (*model.Account, error) {
	output, err := updateAccountRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Account), err
}
