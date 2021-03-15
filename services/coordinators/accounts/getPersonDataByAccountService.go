package accountservicescoordinators

import (
	"errors"

	horeekaaexception "github.com/horeekaa/backend/_errors/repoExceptions"
	horeekaafailure "github.com/horeekaa/backend/_errors/serviceFailures"
	horeekaafailureenums "github.com/horeekaa/backend/_errors/serviceFailures/_enums"

	horeekaafailuretoerror "github.com/horeekaa/backend/_errors/usecaseErrors/_failureToError"
	servicerepodependencies "github.com/horeekaa/backend/dependencies/services/repos"
	model "github.com/horeekaa/backend/model"
	accountservicecoordinatorinterfaces "github.com/horeekaa/backend/services/coordinators/interfaces/accounts"
	servicecoordinatormodels "github.com/horeekaa/backend/services/coordinators/models"
	databaseservicerepointerfaces "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
)

type getPersonDataFromAccountService struct {
	personService                            databaseservicerepointerfaces.PersonService
	accountService                           databaseservicerepointerfaces.AccountService
	getPersonDataFromAccountUsecaseComponent accountservicecoordinatorinterfaces.GetPersonDataFromAccountUsecaseComponent
}

func NewGetPersonDataFromAccount(GetPersonDataFromAccountUsecaseComponent accountservicecoordinatorinterfaces.GetPersonDataFromAccountUsecaseComponent) (accountservicecoordinatorinterfaces.GetPersonDataFromAccountService, error) {
	personService, _ := servicerepodependencies.InitializePersonService()
	accountService, _ := servicerepodependencies.InitializeAccountService()

	return &getPersonDataFromAccountService{
		personService:                            personService,
		accountService:                           accountService,
		getPersonDataFromAccountUsecaseComponent: GetPersonDataFromAccountUsecaseComponent,
	}, nil
}

func (getPrsnData *getPersonDataFromAccountService) preExecute(input model.Account) (*model.Account, error) {
	if &input.ID == nil {
		return &model.Account{}, horeekaafailure.NewFailureObject(
			horeekaafailureenums.AccountIDNeededToRetrievePersonData,
			"/getPersonDataFromAccount",
			horeekaaexception.NewExceptionObject(
				horeekaafailureenums.AccountIDNeededToRetrievePersonData,
				"/getPersonDataFromAccount",
				errors.New(horeekaafailureenums.AccountIDNeededToRetrievePersonData),
			),
		)
	}
	return getPrsnData.getPersonDataFromAccountUsecaseComponent.Validation(input)
}

func (getPrsnData *getPersonDataFromAccountService) Execute(input model.Account) (*servicecoordinatormodels.GetPersonDataByAccountOutput, error) {
	preExecuteOutput, err := getPrsnData.preExecute(input)
	if err != nil {
		return nil, horeekaafailuretoerror.ConvertFailure(
			"/getPersonDataFromAccount",
			&err,
		)
	}
	accountChannel, errChannel := getPrsnData.accountService.FindByID((*preExecuteOutput).ID, &databaseserviceoperations.ServiceOptions{})
	account := &model.Account{}
	select {
	case account = <-accountChannel:
		break
	case err := <-errChannel:
		return nil, horeekaafailuretoerror.ConvertFailure(
			"/getPersonDataFromAccount",
			&err,
		)
	}

	prsonChannel, errChannel := getPrsnData.personService.FindByID((*account).Person.ID, &databaseserviceoperations.ServiceOptions{})
	select {
	case person := <-prsonChannel:
		return &servicecoordinatormodels.GetPersonDataByAccountOutput{
			Person:  person,
			Account: account,
		}, nil

	case err := <-errChannel:
		return nil, horeekaafailuretoerror.ConvertFailure(
			"/getPersonDataFromAccount",
			&err,
		)
	}
}
