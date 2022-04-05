package accountdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type UpdateAccountTransactionComponent interface {
	PreTransaction(
		updateAccountInput *model.InternalUpdateAccount,
	) (*model.InternalUpdateAccount, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateAccountInput *model.InternalUpdateAccount,
	) (*model.Account, error)
}

type UpdateAccountRepository interface {
	RunTransaction(
		updateAccountInput *model.InternalUpdateAccount,
	) (*model.Account, error)
}
