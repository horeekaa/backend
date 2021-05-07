package accountdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type CreateAccountFromAuthDataTransactionComponent interface {
	PreTransaction(
		createAccFromAuthDataInput accountdomainrepositorytypes.CreateAccountFromAuthDataInput,
	) (accountdomainrepositorytypes.CreateAccountFromAuthDataInput, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createAccFromAuthDataInput accountdomainrepositorytypes.CreateAccountFromAuthDataInput,
	) (*model.Account, error)
}

type CreateAccountFromAuthDataRepository interface {
	RunTransaction(
		createAccFromAuthDataInput accountdomainrepositorytypes.CreateAccountFromAuthDataInput,
	) (*model.Account, error)
}
