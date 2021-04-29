package organizationdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetOrganizationRepository interface {
	Execute(filterFields *model.UpdateOrganization) (*model.Organization, error)
}
