package accountpresentationusecases

import (
	"errors"

	horeekaacoreerror "github.com/horeekaa/backend/core/_errors/usecaseErrors"
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/_errors/usecaseErrors/_failureToError"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getPersonDataFromAccountUsecase struct {
	getPersonDataFromAccountRepository accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository
}

func NewGetPersonDataFromAccountUsecase(
	getPersonDataFromAccountRepository accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
) (accountpresentationusecaseinterfaces.GetPersonDataFromAccountUsecase, error) {
	return &getPersonDataFromAccountUsecase{
		getPersonDataFromAccountRepository,
	}, nil
}

func (getPersonDataFromAccountUsecase *getPersonDataFromAccountUsecase) validation(input model.Account) (*model.Account, error) {
	if &input.ID == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			"getPersonDataFromAccount/",
			404,
			"test",
			errors.New("getPersonDataFromAccount/"),
		)
	}
	return &input, nil
}

func (getPersonDataFromAccountUsecase *getPersonDataFromAccountUsecase) Execute(input model.Account) (*model.Person, error) {
	validatedInput, err := getPersonDataFromAccountUsecase.validation(input)
	if err != nil {
		return nil, err
	}

	result, err := getPersonDataFromAccountUsecase.getPersonDataFromAccountRepository.Execute(validatedInput)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"getPersonDataFromAccount/",
			err,
		)
	}
	return result.Person, nil
}
