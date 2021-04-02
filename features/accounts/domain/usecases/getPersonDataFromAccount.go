package accountpresentationusecases

import (
	horeekaacorebasefailure "github.com/horeekaa/backend/core/_errors/serviceFailures/_base"
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

func (getPersonDataFromAccountUsecase *getPersonDataFromAccountUsecase) Validation(input model.Account) (*model.Account, error) {
	if &input.ID == nil {
		return nil, horeekaacoreerror.NewErrorObject(
			"getPersonDataFromAccount/",
			404,
			"test",
			&horeekaacorebasefailure.Failure{},
		)
	}
	return &input, nil
}

func (getPersonDataFromAccountUsecase *getPersonDataFromAccountUsecase) Execute(input *model.Account) (*model.Person, error) {
	result, err := getPersonDataFromAccountUsecase.getPersonDataFromAccountRepository.Execute(*input)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"getPersonDataFromAccount/",
			&err,
		)
	}
	return result.Person, nil
}
