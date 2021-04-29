package organizationdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type UpdateOrganizationOutput struct {
	PreviousOrganization *model.Organization
	UpdatedOrganization  *model.Organization
}

type GetAllOrganizationInput struct {
	FilterFields  *model.UpdateOrganization
	PaginationOpt *model.PaginationOptionInput
}
