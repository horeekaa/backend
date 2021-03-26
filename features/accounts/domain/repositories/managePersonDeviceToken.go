package accountdomainrepositoryinterfaces

import (
	accountrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	model "github.com/horeekaa/backend/model"
)

type ManagePersonDeviceTokenUsecaseComponent interface {
	Validation(input accountrepositorytypes.ManagePersonDeviceTokenInput) (*accountrepositorytypes.ManagePersonDeviceTokenInput, error)
}

type ManagePersonDeviceTokenRepository interface {
	SetValidation(usecaseComponent ManagePersonDeviceTokenUsecaseComponent) (bool, error)
	Execute(input accountrepositorytypes.ManagePersonDeviceTokenInput) (*model.Person, error)
}
