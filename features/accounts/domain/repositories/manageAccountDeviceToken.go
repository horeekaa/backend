package accountdomainrepositoryinterfaces

import (
	accountrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	model "github.com/horeekaa/backend/model"
)

type ManageAccountDeviceTokenUsecaseComponent interface {
	Validation(input accountrepositorytypes.ManageAccountDeviceTokenInput) (accountrepositorytypes.ManageAccountDeviceTokenInput, error)
}

type ManageAccountDeviceTokenRepository interface {
	SetValidation(usecaseComponent ManageAccountDeviceTokenUsecaseComponent) (bool, error)
	Execute(input accountrepositorytypes.ManageAccountDeviceTokenInput) (*model.Account, error)
}
