package organizationpresentationusecasetypes

import (
	"context"

	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"
	"github.com/horeekaa/backend/model"
)

type CreateOrganizationUsecaseInput struct {
	User               *authenticationcoremodels.AuthUserWrap
	Context            context.Context
	CreateOrganization *model.CreateOrganization
}

type UpdateOrganizationUsecaseInput struct {
	User               *authenticationcoremodels.AuthUserWrap
	Context            context.Context
	UpdateOrganization *model.UpdateOrganization
}

type GetAllOrganizationUsecaseInput struct {
	User          *authenticationcoremodels.AuthUserWrap
	Context       context.Context
	FilterFields  *model.OrganizationFilterFields
	PaginationOps *model.PaginationOptionInput
}
