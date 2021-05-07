package accountdomainrepositoryinterfaces

import (
	"firebase.google.com/go/v4/auth"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
)

type GetUserFromAuthHeaderRepository interface {
	Execute(input accountdomainrepositorytypes.GetUserFromAuthHeaderInput) (*auth.UserRecord, error)
}
