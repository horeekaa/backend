package memberaccessrefdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	memberaccessrefdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type UpdateMemberAccessRefUsecaseComponent interface {
	Validation(
		updateAccountAuthInput *model.UpdateMemberAccessRef,
	) (*model.UpdateMemberAccessRef, error)
}

type UpdateMemberAccessRefTransactionComponent interface {
	SetValidation(usecaseComponent UpdateMemberAccessRefUsecaseComponent) (bool, error)

	PreTransaction(
		updateAccountAuthInput *model.UpdateMemberAccessRef,
	) (*model.UpdateMemberAccessRef, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateMemberAccessRefInput *model.UpdateMemberAccessRef,
	) (*memberaccessrefdomainrepositorytypes.UpdateMemberAccessRefOutput, error)
}

type UpdateMemberAccessRefRepository interface {
	SetValidation(usecaseComponent UpdateMemberAccessRefUsecaseComponent) (bool, error)
	RunTransaction(
		updateMemberAccessRefInput *model.UpdateMemberAccessRef,
	) (*memberaccessrefdomainrepositorytypes.UpdateMemberAccessRefOutput, error)
}
