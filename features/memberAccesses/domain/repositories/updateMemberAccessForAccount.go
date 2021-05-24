package memberaccessdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type UpdateMemberAccessForAccountTransactionComponent interface {
	PreTransaction(
		updateMemberAccessInput *model.UpdateMemberAccess,
	) (*model.UpdateMemberAccess, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateMemberAccessInput *model.UpdateMemberAccess,
	) (*memberaccessdomainrepositorytypes.UpdateMemberAccessOutput, error)
}

type UpdateMemberAccessForAccountRepository interface {
	RunTransaction(
		updateMemberAccessInput *model.UpdateMemberAccess,
	) (*memberaccessdomainrepositorytypes.UpdateMemberAccessOutput, error)
}
