package organizationpresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreateOrganizationUsecaseInput struct {
	AuthHeader         string
	Context            context.Context
	CreateOrganization *model.CreateOrganization
}

type UpdateOrganizationUsecaseInput struct {
	AuthHeader         string
	Context            context.Context
	UpdateOrganization *model.UpdateOrganization
}

type GetAllOrganizationUsecaseInput struct {
	AuthHeader    string
	Context       context.Context
	FilterFields  *model.UpdateOrganization
	PaginationOps *model.PaginationOptionInput
}
