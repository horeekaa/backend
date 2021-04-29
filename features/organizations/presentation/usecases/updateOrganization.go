package organizationpresentationusecaseinterfaces

import (
	organizationpresentationusecasetypes "github.com/horeekaa/backend/features/organizations/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdateOrganizationUsecase interface {
	Execute(input organizationpresentationusecasetypes.UpdateOrganizationUsecaseInput) (*model.Organization, error)
}
