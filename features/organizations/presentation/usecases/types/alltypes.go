package organizationpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreateOrganizationUsecaseInput struct {
	Context            context.Context
	CreateOrganization *model.CreateOrganization
}

type UpdateOrganizationUsecaseInput struct {
	Context            context.Context
	UpdateOrganization *model.UpdateOrganization
}

type GetAllOrganizationUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.OrganizationFilterFields
	PaginationOps *model.PaginationOptionInput
}
