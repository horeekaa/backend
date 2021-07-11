package memberaccessdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessUsecaseComponent interface {
	Validation(
		createMemberAccessInput *model.InternalCreateMemberAccess,
	) (*model.InternalCreateMemberAccess, error)
}

type CreateMemberAccessTransactionComponent interface {
	SetValidation(usecaseComponent CreateMemberAccessUsecaseComponent) (bool, error)

	PreTransaction(
		createMemberAccessInput *model.InternalCreateMemberAccess,
	) (*model.InternalCreateMemberAccess, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createMemberAccessInput *model.InternalCreateMemberAccess,
	) (*model.MemberAccess, error)
}

type CreateMemberAccessRepository interface {
	SetValidation(usecaseComponent CreateMemberAccessUsecaseComponent) (bool, error)
	RunTransaction(
		createMemberAccessInput *model.InternalCreateMemberAccess,
	) (*model.MemberAccess, error)
}
