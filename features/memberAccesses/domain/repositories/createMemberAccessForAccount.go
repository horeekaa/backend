package memberaccessdomainrepositoryinterfaces

import (
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessForAccountUsecaseComponent interface {
	Validation(input memberaccessdomainrepositorytypes.CreateMemberAccessForAccountInput) (memberaccessdomainrepositorytypes.CreateMemberAccessForAccountInput, error)
}

type CreateMemberAccessForAccountRepository interface {
	SetValidation(usecaseComponent CreateMemberAccessForAccountUsecaseComponent) (bool, error)
	Execute(input memberaccessdomainrepositorytypes.CreateMemberAccessForAccountInput) (*model.MemberAccess, error)
}
