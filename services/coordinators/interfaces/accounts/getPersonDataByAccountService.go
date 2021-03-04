package accountservicecoordinatorinterfaces

import (
	model "github.com/horeekaa/backend/model"
	servicecoordinatormodels "github.com/horeekaa/backend/services/coordinators/models"
)

type GetPersonDataFromAccountUsecaseComponent interface {
	Validation(input model.Account) (model.Account, error)
}

type GetPersonDataFromAccountService interface {
	Execute(input model.Account) (*servicecoordinatormodels.GetPersonDataByAccountOutput, error)
}
