package accountdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type ManageAccountAuthenticationUsecaseComponent interface {
	Validation(
		manageAccountAuthInput accountdomainrepositorytypes.ManageAccountAuthenticationInput,
	) (accountdomainrepositorytypes.ManageAccountAuthenticationInput, error)
}

type ManageAccountAuthenticationTransactionComponent interface {
	SetValidation(usecaseComponent ManageAccountAuthenticationUsecaseComponent) (bool, error)

	PreTransaction(
		manageAccountAuthInput accountdomainrepositorytypes.ManageAccountAuthenticationInput,
	) (accountdomainrepositorytypes.ManageAccountAuthenticationInput, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		manageAccountAuthInput accountdomainrepositorytypes.ManageAccountAuthenticationInput,
	) (*model.Account, error)
}

type ManageAccountAuthenticationRepository interface {
	SetValidation(usecaseComponent ManageAccountAuthenticationUsecaseComponent) (bool, error)
	RunTransaction(
		manageAccountAuthInput accountdomainrepositorytypes.ManageAccountAuthenticationInput,
	) (*model.Account, error)
}
