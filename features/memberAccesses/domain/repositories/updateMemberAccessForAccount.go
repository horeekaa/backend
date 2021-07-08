package memberaccessdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type UpdateMemberAccessForAccountTransactionComponent interface {
	PreTransaction(
		updateMemberAccessInput *model.InternalUpdateMemberAccess,
	) (*model.InternalUpdateMemberAccess, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateMemberAccessInput *model.InternalUpdateMemberAccess,
	) (*model.MemberAccess, error)
}

type UpdateMemberAccessForAccountRepository interface {
	RunTransaction(
		updateMemberAccessInput *model.InternalUpdateMemberAccess,
	) (*model.MemberAccess, error)
}
