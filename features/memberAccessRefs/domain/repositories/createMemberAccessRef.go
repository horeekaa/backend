package memberaccessrefdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessRefUsecaseComponent interface {
	Validation(
		createMemberAccessRefInput *model.InternalCreateMemberAccessRef,
	) (*model.InternalCreateMemberAccessRef, error)
}

type CreateMemberAccessRefTransactionComponent interface {
	SetValidation(usecaseComponent CreateMemberAccessRefUsecaseComponent) (bool, error)

	PreTransaction(
		createMemberAccessRefInput *model.InternalCreateMemberAccessRef,
	) (*model.InternalCreateMemberAccessRef, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createMemberAccessRefInput *model.InternalCreateMemberAccessRef,
	) (*model.MemberAccessRef, error)
}

type CreateMemberAccessRefRepository interface {
	SetValidation(usecaseComponent CreateMemberAccessRefUsecaseComponent) (bool, error)
	RunTransaction(
		createMemberAccessRefInput *model.InternalCreateMemberAccessRef,
	) (*model.MemberAccessRef, error)
}
