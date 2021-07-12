package memberaccessdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ApproveUpdateMemberAccessUsecaseComponent interface {
	Validation(
		updateMemberAccessInput *model.InternalUpdateMemberAccess,
	) (*model.InternalUpdateMemberAccess, error)
}

type ApproveUpdateMemberAccessTransactionComponent interface {
	SetValidation(usecaseComponent ApproveUpdateMemberAccessUsecaseComponent) (bool, error)

	PreTransaction(
		updateMemberAccessInput *model.InternalUpdateMemberAccess,
	) (*model.InternalUpdateMemberAccess, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateMemberAccessInput *model.InternalUpdateMemberAccess,
	) (*model.MemberAccess, error)
}

type ApproveUpdateMemberAccessRepository interface {
	SetValidation(usecaseComponent ApproveUpdateMemberAccessUsecaseComponent) (bool, error)
	RunTransaction(
		updateMemberAccessInput *model.InternalUpdateMemberAccess,
	) (*model.MemberAccess, error)
}
