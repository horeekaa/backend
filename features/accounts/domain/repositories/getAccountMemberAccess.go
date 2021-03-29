package accountdomainrepositoryinterfaces

import (
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAccountMemberAccessUsecaseComponent interface {
	Validation(input accountdomainrepositorytypes.GetAccountMemberAccessInput) (*accountdomainrepositorytypes.GetAccountMemberAccessInput, error)
}

type GetAccountMemberAccessRepository interface {
	SetValidation(usecaseComponent GetAccountMemberAccessUsecaseComponent) (bool, error)
	Execute(input accountdomainrepositorytypes.GetAccountMemberAccessInput) (*model.MemberAccess, error)
}
