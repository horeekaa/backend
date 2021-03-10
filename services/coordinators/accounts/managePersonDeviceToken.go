package accountservicescoordinators

import (
	horeekaafailuretoerror "github.com/horeekaa/backend/_errors/usecaseErrors/_failureToError"
	servicerepodependencies "github.com/horeekaa/backend/dependencies/services/repos"
	"github.com/horeekaa/backend/model"
	accountservicecoordinatorinterfaces "github.com/horeekaa/backend/services/coordinators/interfaces/accounts"
	servicecoordinatorinterfaces "github.com/horeekaa/backend/services/coordinators/interfaces/accounts"
	servicecoordinatormodels "github.com/horeekaa/backend/services/coordinators/models"
	databaseservicerepointerfaces "github.com/horeekaa/backend/services/database/interfaces/repos"
	databaseserviceoperations "github.com/horeekaa/backend/services/database/operations"
	"github.com/thoas/go-funk"
)

type managePersonDeviceTokenService struct {
	personService                           databaseservicerepointerfaces.PersonService
	managePersonDeviceTokenUsecaseComponent servicecoordinatorinterfaces.ManagePersonDeviceTokenUsecaseComponent
	managePersonDeviceTokenType             string
}

func NewManagePersonDeviceToken(managePersonDeviceTokenComponent accountservicecoordinatorinterfaces.ManagePersonDeviceTokenUsecaseComponent) (accountservicecoordinatorinterfaces.ManagePersonDeviceTokenService, error) {
	personService, _ := servicerepodependencies.InitializePersonService()

	return &managePersonDeviceTokenService{
		personService:                           personService,
		managePersonDeviceTokenUsecaseComponent: managePersonDeviceTokenComponent,
	}, nil
}

func (mgsAccDevToken *managePersonDeviceTokenService) preExecute(input servicecoordinatormodels.ManagePersonDeviceTokenInput) (*servicecoordinatormodels.ManagePersonDeviceTokenInput, error) {
	return mgsAccDevToken.managePersonDeviceTokenUsecaseComponent.Validation(input)
}

func (mgsPrsDevToken *managePersonDeviceTokenService) Execute(input servicecoordinatormodels.ManagePersonDeviceTokenInput) (*model.Person, error) {
	_, err := mgsPrsDevToken.preExecute(input)
	if err != nil {
		return nil, horeekaafailuretoerror.ConvertFailure(
			"/getPersonDataByAccount",
			&err,
		)
	}
	prsnChannel, errChannel := mgsPrsDevToken.personService.FindByID(input.Person.ID, &databaseserviceoperations.ServiceOptions{})

	select {
	case person := <-prsnChannel:
		switch input.ManagePersonDeviceTokenAction {
		case servicecoordinatormodels.ManagePersonDeviceTokenActionInsert:
			if !funk.Contains(person.DeviceTokens, input.DeviceToken) {
				person.DeviceTokens = append(person.DeviceTokens, &input.DeviceToken)
			}
		case servicecoordinatormodels.ManagePersonDeviceTokenActionRemove:
			index := funk.IndexOf(person.DeviceTokens, input.DeviceToken)
			person.DeviceTokens = append(person.DeviceTokens[:index], person.DeviceTokens[index+1:]...)
		}

		prsnChannel, errChannel := mgsPrsDevToken.personService.Update(
			person.ID,
			&model.UpdatePerson{
				DeviceTokens: person.DeviceTokens,
			},
			&databaseserviceoperations.ServiceOptions{},
		)
		select {
		case updatedPerson := <-prsnChannel:
			return updatedPerson, nil

		case err := <-errChannel:
			return nil, horeekaafailuretoerror.ConvertFailure(
				"/getPersonDataFromAccount",
				&err,
			)
		}

	case err := <-errChannel:
		return nil, horeekaafailuretoerror.ConvertFailure(
			"/getPersonDataFromAccount",
			&err,
		)
	}
}
