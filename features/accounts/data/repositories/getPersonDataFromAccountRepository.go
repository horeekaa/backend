package accountdomainrepositories

import (
	"errors"

	horeekaacorefailure "github.com/horeekaa/backend/core/_errors/serviceFailures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/_errors/serviceFailures/_enums"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"

	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/_errors/usecaseErrors/_failureToError"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	model "github.com/horeekaa/backend/model"
)

type getPersonDataFromAccountRepository struct {
	personDataSource                         databaseaccountdatasourceinterfaces.PersonDataSource
	accountDataSource                        databaseaccountdatasourceinterfaces.AccountDataSource
	getPersonDataFromAccountUsecaseComponent accountdomainrepositoryinterfaces.GetPersonDataFromAccountUsecaseComponent
}

func NewGetPersonDataFromAccountRepository(
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
) (accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository, error) {
	return &getPersonDataFromAccountRepository{
		personDataSource:  personDataSource,
		accountDataSource: accountDataSource,
	}, nil
}

func (getPrsnData *getPersonDataFromAccountRepository) SetValidation(
	usecaseComponent accountdomainrepositoryinterfaces.GetPersonDataFromAccountUsecaseComponent,
) (bool, error) {
	getPrsnData.getPersonDataFromAccountUsecaseComponent = usecaseComponent
	return true, nil
}

func (getPrsnData *getPersonDataFromAccountRepository) preExecute(input model.Account) (*model.Account, error) {
	if &input.ID == nil {
		return &model.Account{}, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.AccountIDNeededToRetrievePersonData,
			"/getPersonDataFromAccount",
			errors.New(horeekaacorefailureenums.AccountIDNeededToRetrievePersonData),
		)
	}
	if getPrsnData.getPersonDataFromAccountUsecaseComponent == nil {
		return &input, nil
	}
	return getPrsnData.getPersonDataFromAccountUsecaseComponent.Validation(input)
}

func (getPrsnData *getPersonDataFromAccountRepository) Execute(input model.Account) (*accountrepositorytypes.GetPersonDataByAccountOutput, error) {
	preExecuteOutput, err := getPrsnData.preExecute(input)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getPersonDataFromAccount",
			err,
		)
	}
	account, err := getPrsnData.accountDataSource.GetMongoDataSource().FindByID((*preExecuteOutput).ID, &mongodbcoretypes.OperationOptions{})
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getPersonDataFromAccount",
			err,
		)
	}

	person, err := getPrsnData.personDataSource.GetMongoDataSource().FindByID((*account).Person.ID, &mongodbcoretypes.OperationOptions{})
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getPersonDataFromAccount",
			err,
		)
	}
	return &accountrepositorytypes.GetPersonDataByAccountOutput{
		Person:  person,
		Account: account,
	}, nil
}
