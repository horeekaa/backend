package accountdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
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
	mongoDBTransaction.SetTransaction(
		createAccountFromAuthDataTransactionComponent,
		"CreateAccountFromAuthDataRepository",
	)
	return &createAccountFromAuthDataRepository{
		createAccountFromAuthDataTransactionComponent,
		mongoDBTransaction,
	}, nil
}

func (createAccFromAuthData *createAccountFromAuthDataRepository) RunTransaction(
	createAccFromAuthDataInput accountdomainrepositorytypes.CreateAccountFromAuthDataInput,
) (*model.Account, error) {
	output, err := createAccFromAuthData.mongoDBTransaction.RunTransaction(createAccFromAuthDataInput)
	return (output).(*model.Account), err
}
