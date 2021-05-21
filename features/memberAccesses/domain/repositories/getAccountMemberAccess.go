package memberaccessdomainrepositoryinterfaces

import (
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAccountMemberAccessUsecaseComponent interface {
	Validation(input memberaccessdomainrepositorytypes.GetAccountMemberAccessInput) (memberaccessdomainrepositorytypes.GetAccountMemberAccessInput, error)
}

type GetAccountMemberAccessRepository interface {
	SetValidation(usecaseComponent GetAccountMemberAccessUsecaseComponent) (bool, error)
	Execute(input memberaccessdomainrepositorytypes.GetAccountMemberAccessInput) (*model.MemberAccess, error)
}
