package accountservicecoordinatorinterfaces

import (
	model "github.com/horeekaa/backend/model"
	servicecoordinatormodels "github.com/horeekaa/backend/services/coordinators/models"
)

type ManagePersonDeviceTokenUsecaseComponent interface {
	Validation(input servicecoordinatormodels.ManagePersonDeviceTokenInput) (*servicecoordinatormodels.ManagePersonDeviceTokenInput, error)
}

type ManagePersonDeviceTokenService interface {
	Execute(input servicecoordinatormodels.ManagePersonDeviceTokenInput) (*model.Person, error)
}
