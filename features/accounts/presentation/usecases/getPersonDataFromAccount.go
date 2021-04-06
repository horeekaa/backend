package accountpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetPersonDataFromAccountUsecase interface {
	Execute(input model.Account) (*model.Person, error)
}
