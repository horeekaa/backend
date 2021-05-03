package accountdomainrepositories

import (
	"errors"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/serviceFailures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/serviceFailures/_enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/serviceFailures/_exceptionToFailure"

	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	model "github.com/horeekaa/backend/model"
)

type getPersonDataFromAccountRepository struct {
	personDataSource                         databaseaccountdatasourceinterfaces.PersonDataSource
	getPersonDataFromAccountUsecaseComponent accountdomainrepositoryinterfaces.GetPersonDataFromAccountUsecaseComponent
}

func NewGetPersonDataFromAccountRepository(
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
) (accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository, error) {
	return &getPersonDataFromAccountRepository{
		personDataSource: personDataSource,
	}, nil
}

func (getPrsnData *getPersonDataFromAccountRepository) SetValidation(
	usecaseComponent accountdomainrepositoryinterfaces.GetPersonDataFromAccountUsecaseComponent,
) (bool, error) {
	getPrsnData.getPersonDataFromAccountUsecaseComponent = usecaseComponent
	return true, nil
}

func (getPrsnData *getPersonDataFromAccountRepository) preExecute(input *model.Account) (*model.Account, error) {
	if &input.ID == nil {
		return nil, horeekaacorefailure.NewFailureObject(
			horeekaacorefailureenums.AccountIDNeededToRetrievePersonData,
			"/getPersonDataFromAccount",
			errors.New(horeekaacorefailureenums.AccountIDNeededToRetrievePersonData),
		)
	}
	if getPrsnData.getPersonDataFromAccountUsecaseComponent == nil {
		return input, nil
	}
	return getPrsnData.getPersonDataFromAccountUsecaseComponent.Validation(input)
}

func (getPrsnData *getPersonDataFromAccountRepository) Execute(input *model.Account) (*model.Person, error) {
	preExecuteOutput, err := getPrsnData.preExecute(input)
	if err != nil {
		return nil, err
	}

	person, err := getPrsnData.personDataSource.GetMongoDataSource().FindByID(preExecuteOutput.Person.ID, &mongodbcoretypes.OperationOptions{})
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getPersonDataFromAccount",
			err,
		)
	}
	return person, nil
}
