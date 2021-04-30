package accountdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type manageAccountAuthenticationRepository struct {
	manageAccountAuthenticationTransactionComponent accountdomainrepositoryinterfaces.ManageAccountAuthenticationTransactionComponent
	mongoDBTransaction                              mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewManageAccountAuthenticationRepository(
	manageAccountAuthenticationTransactionComponent accountdomainrepositoryinterfaces.ManageAccountAuthenticationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository, error) {
	mongoDBTransaction.SetTransaction(
		manageAccountAuthenticationTransactionComponent,
		"ManageAccountRepository",
	)
	return &manageAccountAuthenticationRepository{
		manageAccountAuthenticationTransactionComponent,
		mongoDBTransaction,
	}, nil
}

func (manageAccAuthRepo *manageAccountAuthenticationRepository) SetValidation(usecaseComponent accountdomainrepositoryinterfaces.ManageAccountAuthenticationUsecaseComponent) (bool, error) {
	manageAccAuthRepo.manageAccountAuthenticationTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (manageAccAuthRepo *manageAccountAuthenticationRepository) RunTransaction(
	manageAccountAuthInput accountdomainrepositorytypes.ManageAccountAuthenticationInput,
) (*model.Account, error) {
	output, err := manageAccAuthRepo.mongoDBTransaction.RunTransaction(manageAccountAuthInput)
	return (output).(*model.Account), err
}
