package organizationdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type GetAllOrganizationInput struct {
	FilterFields  *model.OrganizationFilterFields
	PaginationOpt *model.PaginationOptionInput
}
