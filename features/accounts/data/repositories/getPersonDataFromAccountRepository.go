package accountdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"

	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	model "github.com/horeekaa/backend/model"
)

type getPersonDataFromAccountRepository struct {
	personDataSource                         databaseaccountdatasourceinterfaces.PersonDataSource
	getPersonDataFromAccountUsecaseComponent accountdomainrepositoryinterfaces.GetPersonDataFromAccountUsecaseComponent
	pathIdentity                             string
}

func NewGetPersonDataFromAccountRepository(
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
) (accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository, error) {
	return &getPersonDataFromAccountRepository{
		personDataSource: personDataSource,
		pathIdentity:     "GetPersonDataFromAccount",
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
			getPrsnData.pathIdentity,
			nil,
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
			getPrsnData.pathIdentity,
			err,
		)
	}
	return person, nil
}
