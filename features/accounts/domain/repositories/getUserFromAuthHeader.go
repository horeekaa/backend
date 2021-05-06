package accountdomainrepositoryinterfaces

import (
	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
)

type GetUserFromAuthHeaderRepository interface {
	Execute(input accountdomainrepositorytypes.GetUserFromAuthHeaderInput) (*authenticationcoremodels.AuthUserWrap, error)
}
