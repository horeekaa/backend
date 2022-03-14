package organizationpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationpresentationusecaseinterfaces "github.com/horeekaa/backend/features/organizations/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getOrganizationUsecase struct {
	getOrganizationRepository organizationdomainrepositoryinterfaces.GetOrganizationRepository
	pathIdentity              string
}

func NewGetOrganizationUsecase(
	getOrganizationRepository organizationdomainrepositoryinterfaces.GetOrganizationRepository,
) (organizationpresentationusecaseinterfaces.GetOrganizationUsecase, error) {
	return &getOrganizationUsecase{
		getOrganizationRepository,
		"GetOrganizationUsecase",
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
			getOrgUcase.pathIdentity,
			err,
		)
	}
	return organization, nil
}
