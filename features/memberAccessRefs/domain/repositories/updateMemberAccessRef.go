package memberaccessrefdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type UpdateMemberAccessRefUsecaseComponent interface {
	Validation(
		updateAccountAuthInput *model.InternalUpdateMemberAccessRef,
	) (*model.InternalUpdateMemberAccessRef, error)
}

type UpdateMemberAccessRefTransactionComponent interface {
	SetValidation(usecaseComponent UpdateMemberAccessRefUsecaseComponent) (bool, error)

	PreTransaction(
		updateAccountAuthInput *model.InternalUpdateMemberAccessRef,
	) (*model.InternalUpdateMemberAccessRef, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateMemberAccessRefInput *model.InternalUpdateMemberAccessRef,
	) (*model.MemberAccessRef, error)
}

type UpdateMemberAccessRefRepository interface {
	SetValidation(usecaseComponent UpdateMemberAccessRefUsecaseComponent) (bool, error)
	RunTransaction(
		updateMemberAccessRefInput *model.InternalUpdateMemberAccessRef,
	) (*model.MemberAccessRef, error)
}
