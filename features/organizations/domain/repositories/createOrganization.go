package organizationdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type CreateOrganizationUsecaseComponent interface {
	Validation(input *model.InternalCreateOrganization) (*model.InternalCreateOrganization, error)
}

type CreateOrganizationRepository interface {
	SetValidation(usecaseComponent CreateOrganizationUsecaseComponent) (bool, error)
	Execute(input *model.InternalCreateOrganization) (*model.Organization, error)
}
