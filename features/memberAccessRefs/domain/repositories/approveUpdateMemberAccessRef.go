package memberaccessrefdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ApproveUpdateMemberAccessRefUsecaseComponent interface {
	Validation(
		updateMemberAccessRefInput *model.InternalUpdateMemberAccessRef,
	) (*model.InternalUpdateMemberAccessRef, error)
}

type ApproveUpdateMemberAccessRefTransactionComponent interface {
	SetValidation(usecaseComponent ApproveUpdateMemberAccessRefUsecaseComponent) (bool, error)

	PreTransaction(
		updateMemberAccessRefInput *model.InternalUpdateMemberAccessRef,
	) (*model.InternalUpdateMemberAccessRef, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateMemberAccessRefInput *model.InternalUpdateMemberAccessRef,
	) (*model.MemberAccessRef, error)
}

type ApproveUpdateMemberAccessRefRepository interface {
	SetValidation(usecaseComponent ApproveUpdateMemberAccessRefUsecaseComponent) (bool, error)
	RunTransaction(
		updateMemberAccessRefInput *model.InternalUpdateMemberAccessRef,
	) (*model.MemberAccessRef, error)
}
