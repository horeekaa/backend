package organizationpresentationusecaseinterfaces

import (
	organizationpresentationusecasetypes "github.com/horeekaa/backend/features/organizations/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreateOrganizationUsecase interface {
	Execute(input organizationpresentationusecasetypes.CreateOrganizationUsecaseInput) (*model.Organization, error)
}
