package organizationpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetOrganizationUsecase interface {
	Execute(input *model.OrganizationFilterFields) (*model.Organization, error)
}
