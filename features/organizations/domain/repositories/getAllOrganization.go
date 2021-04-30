package organizationdomainrepositoryinterfaces

import (
	organizationdomainrepositorytypes "github.com/horeekaa/backend/features/organizations/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllOrganizationRepository interface {
	Execute(filterFields organizationdomainrepositorytypes.GetAllOrganizationInput) ([]*model.Organization, error)
}
