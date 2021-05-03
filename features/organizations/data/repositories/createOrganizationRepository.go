package organizationdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/serviceFailures/_exceptionToFailure"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createOrganizationRepository struct {
	organizationDataSource             databaseorganizationdatasourceinterfaces.OrganizationDataSource
	createOrganizationUsecaseComponent organizationdomainrepositoryinterfaces.CreateOrganizationUsecaseComponent
}

func NewCreateOrganizationRepository(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
) (organizationdomainrepositoryinterfaces.CreateOrganizationRepository, error) {
	return &createOrganizationRepository{
		organizationDataSource: organizationDataSource,
	}, nil
}

func (createOrgRepo *createOrganizationRepository) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.CreateOrganizationUsecaseComponent,
) (bool, error) {
	createOrgRepo.createOrganizationUsecaseComponent = usecaseComponent
	return true, nil
}

func (createOrgRepo *createOrganizationRepository) preExecute(
	input *model.CreateOrganization,
) (*model.CreateOrganization, error) {
	if createOrgRepo.createOrganizationUsecaseComponent == nil {
		return input, nil
	}
	return createOrgRepo.createOrganizationUsecaseComponent.Validation(input)
}

func (createOrgRepo *createOrganizationRepository) Execute(
	input *model.CreateOrganization,
) (*model.Organization, error) {
	validatedInput, err := createOrgRepo.preExecute(input)
	if err != nil {
		return nil, err
	}

	newOrganization, err := createOrgRepo.organizationDataSource.GetMongoDataSource().Create(
		validatedInput,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createOrganization",
			err,
		)
	}
	return newOrganization, nil
}
