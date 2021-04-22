package organizationdomainrepositorytypes

import "github.com/horeekaa/backend/model"

type UpdateOrganizationOutput struct {
	PreviousOrganization *model.Organization
	UpdatedOrganization  *model.Organization
}
