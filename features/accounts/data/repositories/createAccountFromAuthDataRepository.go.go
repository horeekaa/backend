package accountdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type createAccountFromAuthDataRepository struct {
	createAccountFromAuthDataTransactionComponent accountdomainrepositoryinterfaces.CreateAccountFromAuthDataTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateAccountFromAuthDataRepository(
	createAccountFromAuthDataTransactionComponent accountdomainrepositoryinterfaces.CreateAccountFromAuthDataTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (accountdomainrepositoryinterfaces.CreateAccountFromAuthDataRepository, error) {
	createAccountFromAuthRepo := &createAccountFromAuthDataRepository{
		createAccountFromAuthDataTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createAccountFromAuthRepo,
		"CreateAccountFromAuthDataRepository",
	)
	return createAccountFromAuthRepo, nil
}

func (createAccFromAuthData *createAccountFromAuthDataRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createAccFromAuthData.createAccountFromAuthDataTransactionComponent.PreTransaction(
		input.(accountdomainrepositorytypes.CreateAccountFromAuthDataInput),
	)
}

func (createAccFromAuthData *createAccountFromAuthDataRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return createAccFromAuthData.createAccountFromAuthDataTransactionComponent.TransactionBody(
		operationOption,
		input.(accountdomainrepositorytypes.CreateAccountFromAuthDataInput),
	)
}

func (createAccFromAuthData *createAccountFromAuthDataRepository) RunTransaction(
	createAccFromAuthDataInput accountdomainrepositorytypes.CreateAccountFromAuthDataInput,
) (*model.Account, error) {
	output, err := createAccFromAuthData.mongoDBTransaction.RunTransaction(createAccFromAuthDataInput)
	return (output).(*model.Account), err
}
