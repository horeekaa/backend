package organizationdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type CreateOrganizationUsecaseComponent interface {
	Validation(input *model.CreateOrganization) (*model.CreateOrganization, error)
}

type CreateOrganizationRepository interface {
	SetValidation(usecaseComponent CreateOrganizationUsecaseComponent) (bool, error)
	Execute(input *model.CreateOrganization) (*model.Organization, error)
}
