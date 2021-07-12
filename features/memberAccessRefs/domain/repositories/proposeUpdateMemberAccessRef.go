package memberaccessrefdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ProposeUpdateMemberAccessRefUsecaseComponent interface {
	Validation(
		updateMemberAccessRefInput *model.InternalUpdateMemberAccessRef,
	) (*model.InternalUpdateMemberAccessRef, error)
}

type ProposeUpdateMemberAccessRefTransactionComponent interface {
	SetValidation(usecaseComponent ProposeUpdateMemberAccessRefUsecaseComponent) (bool, error)

	PreTransaction(
		updateMemberAccessRefInput *model.InternalUpdateMemberAccessRef,
	) (*model.InternalUpdateMemberAccessRef, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateMemberAccessRefInput *model.InternalUpdateMemberAccessRef,
	) (*model.MemberAccessRef, error)
}

type ProposeUpdateMemberAccessRefRepository interface {
	SetValidation(usecaseComponent ProposeUpdateMemberAccessRefUsecaseComponent) (bool, error)
	RunTransaction(
		updateMemberAccessRefInput *model.InternalUpdateMemberAccessRef,
	) (*model.MemberAccessRef, error)
}
