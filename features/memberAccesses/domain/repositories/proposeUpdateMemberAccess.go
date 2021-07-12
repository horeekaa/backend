package memberaccessdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ProposeUpdateMemberAccessUsecaseComponent interface {
	Validation(
		updateMemberAccessInput *model.InternalUpdateMemberAccess,
	) (*model.InternalUpdateMemberAccess, error)
}

type ProposeUpdateMemberAccessTransactionComponent interface {
	SetValidation(usecaseComponent ProposeUpdateMemberAccessUsecaseComponent) (bool, error)

	PreTransaction(
		updateMemberAccessInput *model.InternalUpdateMemberAccess,
	) (*model.InternalUpdateMemberAccess, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateMemberAccessInput *model.InternalUpdateMemberAccess,
	) (*model.MemberAccess, error)
}

type ProposeUpdateMemberAccessRepository interface {
	SetValidation(usecaseComponent ProposeUpdateMemberAccessUsecaseComponent) (bool, error)
	RunTransaction(
		updateMemberAccessInput *model.InternalUpdateMemberAccess,
	) (*model.MemberAccess, error)
}
