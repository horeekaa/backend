package accountdomainrepositoryinterfaces

import (
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessForAccountUsecaseComponent interface {
	Validation(input accountdomainrepositorytypes.CreateMemberAccessForAccountInput) (accountdomainrepositorytypes.CreateMemberAccessForAccountInput, error)
}

type CreateMemberAccessForAccountRepository interface {
	SetValidation(usecaseComponent CreateMemberAccessForAccountUsecaseComponent) (bool, error)
	Execute(input accountdomainrepositorytypes.CreateMemberAccessForAccountInput) (*model.MemberAccess, error)
}
