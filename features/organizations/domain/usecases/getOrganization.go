package organizationpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/_errors/usecaseErrors/_failureToError"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getOrganizationUsecase struct {
	getOrganizationRepository organizationdomainrepositoryinterfaces.GetOrganizationRepository
}

func NewGetOrganizationUsecase(
	getOrganizationRepository organizationdomainrepositoryinterfaces.GetOrganizationRepository,
) (organizationpresentationusecaseinterfaces.GetOrganizationUsecase, error) {
	return &getOrganizationUsecase{
		getOrganizationRepository,
	}, nil
}

func (getOrgUcase *getOrganizationUsecase) validation(
	input *model.OrganizationFilterFields,
) (*model.OrganizationFilterFields, error) {
	return input, nil
}

func (getOrgUcase *getOrganizationUsecase) Execute(
	filterFields *model.OrganizationFilterFields,
) (*model.Organization, error) {
	validatedFilterFields, err := getOrgUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	organization, err := getOrgUcase.getOrganizationRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getOrganization",
			err,
		)
	}
	return organization, nil
}
