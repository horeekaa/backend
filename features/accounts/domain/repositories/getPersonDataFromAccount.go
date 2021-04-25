package accountdomainrepositoryinterfaces

import (
	model "github.com/horeekaa/backend/model"
)

type GetPersonDataFromAccountUsecaseComponent interface {
	Validation(input *model.Account) (*model.Account, error)
}

type GetPersonDataFromAccountRepository interface {
	SetValidation(usecaseComponent GetPersonDataFromAccountUsecaseComponent) (bool, error)
	Execute(input *model.Account) (*model.Person, error)
}
