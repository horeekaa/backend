package accountdomainrepositoryinterfaces

import (
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAccountFromAuthData interface {
	Execute(input accountdomainrepositorytypes.GetAccountFromAuthDataInput) (*model.Account, error)
}
