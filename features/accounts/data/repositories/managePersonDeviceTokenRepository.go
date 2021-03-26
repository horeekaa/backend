package accountdomainrepositories

import (
	horeekaafailuretoerror "github.com/horeekaa/backend/core/_errors/usecaseErrors/_failureToError"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountdomainrepositorytypes "github.com/horeekaa/backend/features/accounts/domain/repositories/types"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type managePersonDeviceTokenRepository struct {
	personDataSource                        databaseaccountdatasourceinterfaces.PersonDataSource
	managePersonDeviceTokenUsecaseComponent accountdomainrepositoryinterfaces.ManagePersonDeviceTokenUsecaseComponent
}

func NewManagePersonDeviceTokenRepository(
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
) (accountdomainrepositoryinterfaces.ManagePersonDeviceTokenRepository, error) {
	return &managePersonDeviceTokenRepository{
		personDataSource: personDataSource,
	}, nil
}

func (mgsAccDevToken *managePersonDeviceTokenRepository) SetValidation(
	usecaseComponent accountdomainrepositoryinterfaces.ManagePersonDeviceTokenUsecaseComponent,
) (bool, error) {
	mgsAccDevToken.managePersonDeviceTokenUsecaseComponent = usecaseComponent
	return true, nil
}

func (mgsAccDevToken *managePersonDeviceTokenRepository) preExecute(input accountdomainrepositorytypes.ManagePersonDeviceTokenInput) (*accountdomainrepositorytypes.ManagePersonDeviceTokenInput, error) {
	return mgsAccDevToken.managePersonDeviceTokenUsecaseComponent.Validation(input)
}

func (mgsPrsDevToken *managePersonDeviceTokenRepository) Execute(input accountdomainrepositorytypes.ManagePersonDeviceTokenInput) (*model.Person, error) {
	_, err := mgsPrsDevToken.preExecute(input)
	if err != nil {
		return nil, horeekaafailuretoerror.ConvertFailure(
			"/getPersonDataByAccount",
			&err,
		)
	}
	person, err := mgsPrsDevToken.personDataSource.GetMongoDataSource().FindByID(input.Person.ID, &mongodbcoretypes.OperationOptions{})
	if err != nil {
		return nil, horeekaafailuretoerror.ConvertFailure(
			"/getPersonDataFromAccount",
			&err,
		)
	}

	switch input.ManagePersonDeviceTokenAction {
	case accountdomainrepositorytypes.ManagePersonDeviceTokenActionInsert:
		if !funk.Contains(person.DeviceTokens, input.DeviceToken) {
			person.DeviceTokens = append(person.DeviceTokens, &input.DeviceToken)
		}
		break

	case accountdomainrepositorytypes.ManagePersonDeviceTokenActionRemove:
		index := funk.IndexOf(person.DeviceTokens, input.DeviceToken)
		person.DeviceTokens = append(person.DeviceTokens[:index], person.DeviceTokens[index+1:]...)
		break
	}

	updatedPerson, err := mgsPrsDevToken.personDataSource.GetMongoDataSource().Update(
		person.ID,
		&model.UpdatePerson{
			DeviceTokens: person.DeviceTokens,
		},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaafailuretoerror.ConvertFailure(
			"/getPersonDataFromAccount",
			&err,
		)
	}
	return updatedPerson, nil
}
