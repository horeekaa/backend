package organizationpresentationusecaseinterfaces

import (
	organizationpresentationusecasetypes "github.com/horeekaa/backend/features/organizations/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllOrganizationUsecase interface {
	Execute(input organizationpresentationusecasetypes.GetAllOrganizationUsecaseInput) ([]*model.Organization, error)
}
