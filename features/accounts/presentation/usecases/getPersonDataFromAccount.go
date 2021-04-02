package accountpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetPersonDataFromAccountUsecase interface {
	Validation(input model.Account) (*model.Account, error)
	Execute(input *model.Account) (*model.Person, error)
}
