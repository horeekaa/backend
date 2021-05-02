package organizationdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetOrganizationRepository interface {
	Execute(filterFields *model.OrganizationFilterFields) (*model.Organization, error)
}
